function __docker_set_find_file()
{
  local found
  local dir=$PWD
  local file=".docker-set"
  while [[ -z "$found" ]] && [[ "$dir" != "" ]]; do
    if [[ -f "$dir/$file" ]]; then
      found="$dir/$file"
      break
    fi
    if [[ "$dir" == "$HOME" ]]; then
      break
    fi
    dir=${dir%/*}
  done
  echo $found
}

function __docker_set_chpwd()
{
  file=$(__docker_set_find_file)
  if [[ -f "$file" ]]; then
    echo "Hello world $file"
  fi
}

export -a chpwd_functions
[[ " ${chpwd_functions[*]} " == *" __docker_set_chpwd "* ]] ||
chpwd_functions+=(__docker_set_chpwd)
 
function __zsh_like_cd()
{
  \typeset __zsh_like_cd_hook
  if
    builtin "$@"
  then
    for __zsh_like_cd_hook in chpwd "${chpwd_functions[@]}"
    do
      if \typeset -f "$__zsh_like_cd_hook" >/dev/null 2>&1
      then "$__zsh_like_cd_hook" || break # finish on first failed hook
      fi
    done
    true
  else
    return $?
  fi
}

[[ -n "${ZSH_VERSION:-}" ]] ||
{
  function cd()    { __zsh_like_cd cd    "$@" ; }
  function popd()  { __zsh_like_cd popd  "$@" ; }
  function pushd() { __zsh_like_cd pushd "$@" ; }
}
