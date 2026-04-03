---
sidebar_position: 3
---

# Trakt

:::info Backup?
You may consider a backup of your server before starting any import. They are not easily reversible, though we do our best to ensure they are accurate and bug free!
:::

**Note:** Your Trakt profile must be public _during_ this process. You are free to private it again once it completes.

1. Provide your Trakt username in the text box.
2. (Optional) [Provide your own api key](#optional-api-key).
3. Press `Start Import`.

This will be a long process, possibly a couple hours depending on how large of a Trakt history you have. If you think it has frozen or isn't working, try checking your server logs to see if it is doing anything.

### (Optional) API Key

This step is optional because Watcharr comes built in with its own API Key that it can use for your import, however, if you encounter any issues where imports cannot start or keep failing and you see any `403` errors in your server logs, providing your own key may fix the problem.

#### Getting an API Key

1. [Visit this link to create a new Trakt API App](https://trakt.tv/oauth/applications/new).
2. Fill out all required fields (`Name` and `Redirect uri`) with any random data.
   1. You can set the `Redirect uri` to any value the input accepts since we don't use it (eg: `http://localhost`)
3. Click `Save App`.
4. You should now be redirected to your new app. Copy the `Client ID` into the `API Key` textbox in Watcharr.
