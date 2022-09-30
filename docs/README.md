# Documentation

To render and view the documentation locally, follow these istructions:

1. Install [Hugo](https://gohugo.io/getting-started/installing/).
2. Download git submodules (Hugo theme):

   ```bash
   git pull --recurse-submodules
   ```

3. Render and host the static website:

   ```bash
   hugo server --minify
   ```

4. Open http://localhost:1313/ in your browser.
