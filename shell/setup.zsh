spd() {
    if [ "$1" = "add" ] && [ "$#" = 1 ]; then
        SPEEDDIAL_ADD_PRINT_COMMAND=1 ./speeddial add $(fc -ln -1)
    elif [ "$1" = "" ]; then
        print -z $(./speeddial)
    else
        ./speeddial $@
    fi

}
