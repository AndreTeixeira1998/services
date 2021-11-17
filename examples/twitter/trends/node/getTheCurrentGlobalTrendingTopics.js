const { TwitterService } = require("micro-js-client/twitter");

// Get the current global trending topics
async function getTheCurrentGlobalTrendingTopics() {
  let twitterService = new TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.trends({});
  console.log(rsp);
}

getTheCurrentGlobalTrendingTopics();
