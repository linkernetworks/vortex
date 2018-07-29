" Specify a directory for plugins
" - For Neovim: ~/.local/share/nvim/plugged
" - Avoid using standard Vim directory names like 'plugin'
call plug#begin('~/.vim/plugged')

" Make sure you use single quotes

" On-demand loading
Plug 'scrooloose/nerdtree', { 'on':  'NERDTreeToggle' }

Plug 'c9s/helper.vim'
Plug 'c9s/treemenu.vim'
Plug 'c9s/hypergit.vim'
Plug 'c9s/vikube.vim'

Plug 'Valloric/YouCompleteMe', { 'do': './install.py --clang-completer --go-completer' }
Plug 'w0rp/ale'

" Go
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
Plug 'nsf/gocode', { 'rtp': 'vim', 'do': '~/.vim/plugged/gocode/vim/symlink.sh' }

" Color Theme
Plug 'dracula/vim', { 'as': 'dracula' }

Plugin 'rstacruz/sparkup'
" Initialize plugin system
call plug#end()

"   nerd tree explorer"{{{
nmap <silent> <leader>e  :NERDTreeToggle<CR>
nmap <silent> <leader>nf :NERDTreeFind<CR>
nmap <silent> <leader>nm :NERDTreeMirror<CR>

cabbr ntf  NERDTreeFind
cabbr ntm  NERDTreeMirror

function! StartUp()
    if 0 == argc()
        NERDTree
    end
endfunction
autocmd VimEnter * call StartUp()

colorscheme dracula

set expandtab

set sidescroll=1
set sidescrolloff=3
set showfulltag showmatch showcmd showmode
set textwidth=0
set winaltkeys=no showtabline=2 hlsearch
set noswapfile

" Trim white space ===========================
fun! TrimWhitespace()
    let l:save = winsaveview()
    %s/\s\+$//e
    call winrestview(l:save)
endfun
command! TrimWhitespace call TrimWhitespace()
noremap <leader>t :call TrimWhitespace()<CR>

" Plugin Settings
" YValloric/YouCompleteMe: YouCompleteMe
let g:ycm_min_num_of_chars_for_completion = 3
" let g:ycm_autoclose_preview_window_after_completion = 1

" w0rp/ale: Asynchronous Lint Engine
let g:ale_linters = { 'Go': ['golint', 'go vet'] }

" plasticboy/vim-markdown
let g:vim_markdown_folding_disabled = 1
let g:vim_markdown_frontmatter = 1 " YMAL Front Matter
let g:vim_markdown_json_frontmatter = 1 " JSON Front Mattetaugroup

" fatih/vim-go
let g:go_fmt_command = "goimports"
let g:go_fmt_fail_silently = 1
let g:go_fmt_autosave = 1
let g:go_highlight_functions = 1
let g:go_highlight_methods = 1
let g:go_highlight_structs = 1
let g:go_highlight_operators = 1
let g:go_highlight_build_constraints = 1
let g:tagbar_type_go = {
    \ 'ctagstype' : 'go',
    \ 'kinds'     : [
        \ 'p:package',
        \ 'i:imports:1',
        \ 'c:constants',
        \ 'v:variables',
        \ 't:types',
        \ 'n:interfaces',
        \ 'w:fields',
        \ 'e:embedded',
        \ 'm:methods',
        \ 'r:constructor',
        \ 'f:functions'
    \ ],
    \ 'sro' : '.',
    \ 'kind2scope' : {
        \ 't' : 'ctype',
        \ 'n' : 'ntype'
    \ },
    \ 'scope2kind' : {
        \ 'ctype' : 't',
        \ 'ntype' : 'n'
    \ },
    \ 'ctagsbin'  : 'gotags',
    \ 'ctagsargs' : '-sort -silent'
\ }


nmap <silent> <leader>b  :GoBuild<CR>
nmap <silent> <leader>gt  :GoTest<CR>
nmap <silent> <leader>gr  :GoReferrers<CR>
nmap <silent> <leader>gf  :GoFillStruct<CR>
