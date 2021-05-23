package handler

import (
	"context"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	file "github.com/micro/services/file/proto"
	"github.com/micro/services/pkg/tenant"
)

type File struct {
	db model.Model
}

func NewFile() *File {
	i := model.ByEquality("project")
	i.Order.Type = model.OrderTypeUnordered

	db := model.New(
		file.Record{},
		&model.Options{
			Key:     "Id",
			Indexes: []model.Index{i},
		},
	)

	return &File{
		db: db,
	}
}

func (e *File) Read(ctx context.Context, req *file.ReadRequest, rsp *file.ReadResponse) error {
	log.Info("Received File.Read request")

	if len(req.Path) == 0 {
		return errors.BadRequest("file.read", "missing file path")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	var files []*file.Record

	project := tenantId + "/" + req.Project

	// read all the files for the project
	err := e.db.Read(model.QueryEquals("project", project), &files)
	if err != nil {
		return err
	}

	// filter the file
	for _, file := range files {
		if file.Path == req.Path {
			// strip the tenant id
			file.Id = strings.TrimPrefix(file.Id, tenantId+"/")
			file.Project = strings.TrimPrefix(file.Project, tenantId+"/")
			rsp.File = file
		}
	}

	return nil
}

func (e *File) Save(ctx context.Context, req *file.SaveRequest, rsp *file.SaveResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	log.Info("Received File.Save request")

	for _, reqFile := range req.Files {
		// prefix the tenant
		reqFile.Id = tenantId + "/" + reqFile.Id
		reqFile.Project = tenantId + "/" + reqFile.Project

		// create the file
		err := e.db.Create(reqFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *File) List(ctx context.Context, req *file.ListRequest, rsp *file.ListResponse) error {
	log.Info("Received File.List request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	// prefix tenant id
	project := tenantId + "/" + req.Project

	var files []*file.Record

	// read all the files for the project
	if err := e.db.Read(model.QueryEquals("project", project), &files); err != nil {
		return err
	}

	// @todo funnily while this is the archetypical
	// query for the KV store interface, it's not supported by the model
	// so we do client side filtering here
	for _, file := range rsp.Files {
		// strip the prefixes
		file.Id = strings.TrimPrefix(file.Id, tenantId+"/")
		file.Project = strings.TrimPrefix(file.Project, tenantId+"/")

		// if requesting all files or path matches
		if req.Path == "" || strings.HasPrefix(file.Path, req.Path) {
			rsp.Files = append(rsp.Files, file)
		}
	}

	return nil
}