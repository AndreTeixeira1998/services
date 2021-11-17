const { NotesService } = require("micro-js-client/notes");

// List all the notes
async function listAllNotes() {
  let notesService = new NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.list({});
  console.log(rsp);
}

listAllNotes();
