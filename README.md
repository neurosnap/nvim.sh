# nvim.sh

neovim plugin directory search from the terminal

```bash
$ curl https://nvim.sh/s/statusline

Name                          Stars  OpenIssues  Updated               Description                                                                  
nvim-lualine/lualine.nvim     1170   6           2022-01-04T17:30:31Z  A blazing fast and easy to configure neovim statusline plugin written in pure lua.
glepnir/galaxyline.nvim       695    65          2021-04-25T09:04:15Z  neovim statusline plugin written in lua                                      
Famiu/feline.nvim             491    8           2021-12-28T01:46:16Z  A minimal, stylish and customizable statusline for Neovim written in Lua     
windwp/windline.nvim          238    1           2022-01-05T11:22:56Z  Animation statusline, floating window statusline. Use lua + luv make some wind
tjdevries/express_line.nvim   171    15          2021-12-01T21:14:32Z  WIP: Statusline written in pure lua. Supports co-routines, functions and jobs.
datwaft/bubbly.nvim           159    11          2021-11-15T03:59:52Z  Bubbly statusline for neovim                                                 
adelarsq/neoline.vim          154    4           2022-01-03T23:41:54Z  Status Line for Neovim focused on beauty and performance ✅                  
tamton-aquib/staline.nvim     121    0           2021-12-27T06:15:19Z  A modern lightweight statusline and bufferline for neovim in lua. Mainly uses unicode symbols for showing info.
ojroques/nvim-hardline        112    1           2021-12-20T17:27:34Z  A simple Neovim statusline written in Lua                                    
NTBBloodbath/galaxyline.nvim  97     11          2022-01-05T07:26:04Z  neovim statusline plugin written in lua                                      
konapun/vacuumline.nvim       14     0           2021-08-11T01:37:04Z  A prebuilt configuration for galaxyline inspired by airline                  
beauwilliams/statusline.lua   95     0           2021-09-28T00:25:30Z  A zero-config minimal statusline for neovim written in lua featuring awesome integrations and blazing speed!
```

## Usage

Help

```bash
curl https://nvim.sh
```

List all plugins

```bash
curl https://nvim.sh/s
```

Search for a plugin

```bash
curl https://nvim.sh/s/statusline
```

List all tags

```bash
curl https://nvim.sh/t
```

Search for plugins based on tag

```bash
curl https://nvim.sh/t/note-taking
```

## Outputs

JSON - all endpoints support `format` query param

```bash
https://nvim.sh/t/sidebar?format=json
```

## Credits

 - https://neovimcraft.com 
 - https://erock.io
