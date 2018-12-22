#!/bin/sh

# copied from https://github.com/mpapis/bash_zsh_support on Dec 22nd, 2018
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

function docker_set_chpwd_function() {
    echo "Hello World $PWD"
}

export -a chpwd_functions                                             # define hooks as an shell array
[[ " ${chpwd_functions[*]} " == *" docker_set_chpwd_function "* ]] || # prevent double addition
chpwd_functions+=(docker_set_chpwd_function)                          # finally add it to the list
