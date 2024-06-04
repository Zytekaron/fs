# fs

This is an over-engineered file server hosted at https://fs.zyte.dev.

Basic usage:
- File names are shown as-is, directory names end with `/`
- Default directory listing behavior:
  - hyperlinks are shown and if a browser is detected
  - otherwise, names are shown in a list separated by `\n`
- `?viewas=` (overrides all defaults)
  - Files:
    - `?viewas=md` will attempt to do a basic markdown rendering of a file in the browser
  - Directories:
    - `?viewas=hl` will send the directory using html and hyperlinks
    - `?viewas=json` will send the directory list as a json array `["file.txt","dir/"]`
    - `?viewas=json_obj` will wrap the array in an object `{"entries":["file.txt","dir/"]}`

# License
**fs** is licensed under the [MIT License](./LICENSE)
