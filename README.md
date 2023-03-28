# Automatic Ai telegram group manager

Automatic Telegram group manager based on [OpenAi](https://openai.com/) GPT-3.5 API.
Very cheap and easy to use. The Bot is stored as [Google Cloud Function](https://cloud.google.com/functions/), so it can be used for free. Themes for generation posts are received by [Google Pub/Sub](https://cloud.google.com/pubsub) and sent by shedule via [Google Cloud Scheduler](https://cloud.google.com/scheduler).

Example of Telegram Channel: [@ai_talking_to_dev](https://t.me/ai_talking_to_dev)
And here is schedule of posts fo this channel:

![cron screenshot](https://github.com/DeryabinSergey/go-ai-poster/blob/media/storage.png)

In task can be a topic for a post from the [dictionary file](dictionary.json). In this case, we randomly get one of the themes by this key in the Dictionary. Or it can be a theme to generate a post directly.
The Dictionary [file](dictionary.json) is stored in [Google Cloud Storage](https://cloud.google.com/storage/), in the repository, it is just an example. You can use your dictionary. Also, [Google Cloud Storage](https://cloud.google.com/storage/) stored messages log for OpenAi model.

Why storage in Google Cloud Storage?

Because this is free for this usage and does not need a hight performance in our case.