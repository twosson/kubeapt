# Developer Dash UI

## Commands

### `npm run dev`

Same as `npm run start`. Starts a server at `localhost:3000` by default.

To set the base API, you can set it through the env var `API_BASE` (i.e. `API_BASE=http://localhost:11111 npm run dev`)

### `npm run build`

Builds production mode of the single page app into the `/build` directory

### :fire_engine: `npm run fix`
This runs both `eslint --fix`  & `stylelint --fix` over the appropriate files so that you don't have to worry about formatting or your css being valid.
