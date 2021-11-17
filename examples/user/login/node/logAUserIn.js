const { UserService } = require("micro-js-client/user");

// Login using username or email. The response will return a new session for successful login,
// 401 in the case of login failure and 500 for any other error
async function logAuserIn() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.login({
    email: "joe@example.com",
    password: "Password1",
  });
  console.log(rsp);
}

logAuserIn();
