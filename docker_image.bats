#!/usr/bin/env bats

@test "Ruby" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | ruby -nle 'puts \$_'"
  echo "status: ${status}"
  echo "output: ${output}"
  [ "$output" = "シェル芸" ]
}

@test "ccze" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | ccze -A"
  [[ "$output" =~ シェル芸 ]]
}

@test "screen" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "screen -v"
  [[ "$output" =~ Screen ]]
}

@test "tmux" {
  run docker container run --rm ${DOCKER_IMAGE} tmux -c "echo シェル芸"
  [ "$output" = "シェル芸" ]
}

@test "ttyrec" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "ttyrec -h"
  [[ "$output" =~ ttyrec ]]
}

@test "TiMidity++" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "timidity -v"
  [[ "$output" =~ TiMidity\+\+ ]]
}

@test "abcMIDI" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "abc2midi -ver"
  [[ "$output" =~ abc2midi ]]
}

@test "R" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | R -q -e 'cat(readLines(\"stdin\"))'"
  [[ "$output" =~ シェル芸 ]]
}

@test "boxes" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | boxes"
  [[ "$output" =~ \/\*\ シェル芸\ \*\/ ]]
}

# /bin/ash は /bin/dash へのエイリアス, /usr/bin/ash は /usr/bin/dash へのエイリアスで、両方とも同じ
# apt install ash ではエイリアスが作成されるのみ
@test "ash" {
  run docker container run --rm ${DOCKER_IMAGE} ash -c "echo シェル芸"
  [ "$output" = シェル芸 ]
}

@test "yash" {
  run docker container run --rm ${DOCKER_IMAGE} yash -c "echo シェル芸"
  [ "$output" = シェル芸 ]
}

@test "jq" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | jq -Rr '.'"
  [ "$output" = シェル芸 ]
}

@test "Vim" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | vim -es +%p +q! /dev/stdin"
  [ "$output" = シェル芸 ]
}

@test "Emacs" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | emacs -Q --batch --insert /dev/stdin --eval='(princ (buffer-string))'"
  [ "$output" = シェル芸 ]
}

@test "Python2" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | python -c 'import sys;print sys.stdin.readline().rstrip()'"
  [ "$output" = シェル芸 ]
}

@test "Python3" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | python3 -c 'import sys;print(sys.stdin.readline().rstrip())'"
  [ "$output" = シェル芸 ]
}

@test "nkf" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | nkf"
  [ "$output" = シェル芸 ]
}

@test "rs" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | grep -o . | rs -T | tr -d ' '"
  [ "$output" = シェル芸 ]
}

@test "pwgen" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "pwgen -h"
  [ $status -eq 1 ]
  [[ "$output" =~ pwgen ]]
}

@test "bc" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo 'print \"シェル芸\n\"' | bc"
  [ "$output" = "シェル芸" ]
}

@test "Perl" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | perl -nle 'print \$_'"
  [ "$output" = "シェル芸" ]
}

@test "toilet" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | toilet"
  [ "${lines[0]}" = '                                          ' ]
  [ "${lines[1]}" = '   ""m                        m  "m       ' ]
  [ "${lines[2]}" = '  mm                           #  #       ' ]
  [ "${lines[3]}" = '    "    m"      mmm""         #  #   #   ' ]
  [ "${lines[4]}" = '       m"          #mm        m"  # m"    ' ]
  [ "${lines[5]}" = '  "mm""         """"  "      m"   #"      ' ]
  [ "${lines[6]}" = '                                          ' ]
  [ "${lines[7]}" = '                                          ' ]
}

@test "figlet" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo ShellGei | figlet"
  echo "lines[0]: '${lines[0]}'"
  [ "${lines[0]}" = " ____  _          _ _  ____      _ " ]
  [ "${lines[1]}" = "/ ___|| |__   ___| | |/ ___| ___(_)" ]
  [ "${lines[2]}" = "\___ \| '_ \ / _ \ | | |  _ / _ \ |" ]
  [ "${lines[3]}" = " ___) | | | |  __/ | | |_| |  __/ |" ]
  [ "${lines[4]}" = "|____/|_| |_|\___|_|_|\____|\___|_|" ]
}

@test "Haskell" {
  run docker container run --rm ${DOCKER_IMAGE} ghc -e 'putStrLn "シェル芸"'
  [ "$output" = "シェル芸" ]
}

@test "Git" {
  run docker container run --rm ${DOCKER_IMAGE} git version
  [[ "$output" =~ "git version" ]]
}

@test "build-essential" {
  run docker container run --rm ${DOCKER_IMAGE} gcc --version
  [[ "${lines[0]}" =~ gcc ]]
}

@test "mecab" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | mecab -Owakati"
  [ "$output" = "シェル 芸 " ]
}

# @test "wget" {
#   run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 > index.txt && python3 -m http.server & sleep 0.5 && wget -q http://localhost:8000/index.txt -O -"
#   [ "$output" = "シェル芸" ]
# }

@test "curl" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 > index.txt && python3 -m http.server & sleep 0.5 && curl -s http://localhost:8000/index.txt"
  echo "output: '$output'"
  [ "$output" = "シェル芸" ]
}

# @test "npm" {
#   run docker container run --rm ${DOCKER_IMAGE} which npm
#   [ "$output" = "/usr/bin/npm" ]
# }

@test "bsdgames" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo '... .... . .-.. .-.. --. . ..  ...-.-' | morse -d"
  [ "$output" = "SHELLGEI" ]
}

@test "fortune" {
  run docker container run --rm ${DOCKER_IMAGE} fortune
  [ $status -eq 0 ]
}

# 2回指定されている
@test "cowsay" {
  run docker container run --rm ${DOCKER_IMAGE} cowsay シェル芸
  [ "${lines[0]}" = ' __________' ]
  [ "${lines[1]}" = '< シェル芸 >' ]
  [ "${lines[2]}" = ' ----------' ]
  [ "${lines[3]}" = '        \   ^__^' ]
  [ "${lines[4]}" = '         \  (oo)\_______' ]
  [ "${lines[5]}" = '            (__)\       )\/\' ]
  [ "${lines[6]}" = '                ||----w |' ]
  [ "${lines[7]}" = '                ||     ||' ]
}

@test "datamash" {
  run docker container run --rm ${DOCKER_IMAGE} datamash --version
  [[ "${lines[0]}" =~ "datamash (GNU datamash)" ]]
}

@test "gawk" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | gawk '{print \$0}'"
  [ "$output" = "シェル芸" ]
}

@test "libxml2-utils" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo '<?xml version=\"1.0\"?><e>ShellGei</e>' | xmllint --xpath '/e/text()' -"
  [ "$output" = "ShellGei" ]
}

@test "zsh" {
  run docker container run --rm ${DOCKER_IMAGE} zsh -c "echo シェル芸"
  [ "$output" = "シェル芸" ]
}

@test "num-utils" {
  run docker container run --rm ${DOCKER_IMAGE} numaverage -h
  [ "${lines[1]}" = "numaverage : A program for finding the average of numbers." ]
}

# 不要では?
@test "apache2-utils" {
  run docker container run --rm ${DOCKER_IMAGE} ab -V
  [[ "${lines[0]}" =~ "ApacheBench" ]]
}

@test "fish" {
  run docker container run --rm ${DOCKER_IMAGE} fish -c "echo シェル芸"
  [ "$output" = "シェル芸" ]
}

@test "lolcat" {
  run docker container run --rm ${DOCKER_IMAGE} lolcat --version
  [[ "${lines[0]}" =~ "lolcat" ]]
}

@test "nyancat" {
  run docker container run --rm ${DOCKER_IMAGE} nyancat -h
  [ "${lines[0]}" = "Terminal Nyancat" ]
}

@test "ImageMagick" {
  run docker container run --rm ${DOCKER_IMAGE} convert -version
  [[ "${lines[0]}" =~ "Version: ImageMagick" ]]
}

@test "moreutils" {
  run docker container run --rm ${DOCKER_IMAGE} errno 1
  [ "$output" = "EPERM 1 許可されていない操作です" ]
}

# strace は docker 上で実行する場合、--cap-add=SYS_PTRACE と --security-opt="seccomp=unconfined" が必要になるため、不要では
@test "strace" {
  run docker container run --rm ${DOCKER_IMAGE} strace -V
  [[ "${lines[0]}" =~ "strace -- version" ]]
}

@test "whiptail" {
  run docker container run --rm ${DOCKER_IMAGE} whiptail -v
  [[ "$output" =~ "whiptail" ]]
}

@test "pandoc" {
  run docker container run --rm ${DOCKER_IMAGE} pandoc -v
  [[ "${lines[0]}" =~ "pandoc" ]]
}

@test "postgresql" {
  run docker container run --rm ${DOCKER_IMAGE} which psql
  [ "$output" = "/usr/bin/psql" ]
}

@test "uconv" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo 30b730a730eb82b8 | xxd -p -r | uconv -f utf-16be -t utf-8"
  [ "$output" = "シェル芸" ]
}

@test "tcsh" {
  run docker container run --rm ${DOCKER_IMAGE} tcsh -c "echo シェル芸"
  [ "$output" = "シェル芸" ]
}

# 不要?
@test "libskk-dev" {
  run docker container run --rm ${DOCKER_IMAGE} stat /usr/lib/x86_64-linux-gnu/libskk.so
  [ "${lines[0]}" = "  File: /usr/lib/x86_64-linux-gnu/libskk.so -> libskk.so.0.0.0" ]
}

@test "kkc" {
  run docker container run --rm ${DOCKER_IMAGE} kkc help
  [ "${lines[1]}" = "  kkc help コマンド" ]
}

@test "morsegen" {
  run docker container run --rm ${DOCKER_IMAGE} morsegen
  [ $status -eq 1 ]
  [[ "${lines[1]}" =~ "Morse Generator." ]]
}

@test "dc" {
  run docker container run --rm ${DOCKER_IMAGE} dc -V
  [[ "${lines[0]}" =~ "dc" ]]
}

@test "telnet" {
  run docker container run --rm ${DOCKER_IMAGE} telnet -h
  [ $status -eq 1 ]
  [ "${lines[0]}" = "telnet: invalid option -- 'h'" ]
}

@test "busybox" {
  run docker container run --rm ${DOCKER_IMAGE} /bin/busybox echo "シェル芸"
  [ "$output" = "シェル芸" ]
}

@test "parallel" {
  run docker container run --rm ${DOCKER_IMAGE} parallel --version
  [[ "${lines[0]}" =~ "GNU parallel" ]]
}

@test "rename" {
  run docker container run --rm ${DOCKER_IMAGE} rename -V
  [[ "${lines[0]}" =~ "/usr/bin/rename" ]]
}

@test "mt" {
  run docker container run --rm ${DOCKER_IMAGE} mt -v
  [[ "${lines[0]}" =~ "mt-st" ]]
}

@test "ffmpeg" {
  run docker container run --rm ${DOCKER_IMAGE} ffmpeg -version
  [[ "${lines[0]}" =~ "ffmpeg version" ]]
}

@test "kakasi" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | nkf -e | kakasi -JH | nkf -w"
  [ "$output" = "シェルげい" ]
}

@test "dateutils" {
  run docker container run --rm ${DOCKER_IMAGE} /usr/bin/dateutils.dtest -V
  [[ "$output" =~ "datetest" ]]
}

@test "fonts-ipafont" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "fc-list | grep ipa | wc -l"
  [ $output -ge 4 ]
}

@test "fonts-vlgothic" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "fc-list | grep vlgothic | wc -l"
  [ $output -ge 2 ]
}

@test "inkscape" {
  run docker container run --rm ${DOCKER_IMAGE} inkscape --version
  [[ "$output" =~ "Inkscape" ]]
}

@test "gnuplot" {
  run docker container run --rm ${DOCKER_IMAGE} gnuplot -V
  [[ "$output" =~ "gnuplot" ]]
}

@test "qrencode" {
  run docker container run --rm ${DOCKER_IMAGE} qrencode -V
  [[ "${lines[0]}" =~ "qrencode version" ]]
}

@test "fonts-nanum" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "fc-list | grep nanum | wc -l"
  [ $output -ge 10 ]
}

@test "fonts-symbola" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "fc-list | grep Symbola | wc -l"
  [ $output -ge 1 ]
}

@test "fonts-noto-color-emoji" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "fc-list | grep NotoColorEmoji | wc -l"
  [ $output -ge 1 ]
}

@test "sl" {
  run docker container run --rm ${DOCKER_IMAGE} which sl
  [ "$output" = /usr/games/sl ]
}

@test "chromium" {
  run docker container run --rm ${DOCKER_IMAGE} chromium-browser --version
  [[ "$output" =~ "Chromium" ]]
}

@test "nginx" {
  run docker container run --rm ${DOCKER_IMAGE} nginx -v
  [[ "$output" =~ "nginx version:" ]]
}

@test "screenfetch" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "screenfetch -V | sed $'s/\033\[[0-9]m//g'"
  [[ "${lines[0]}" =~ "screenFetch - Version" ]]
}

@test "mono-runtime" {
  run docker container run --rm ${DOCKER_IMAGE} mono --version
  [[ "${lines[0]}" =~ "Mono JIT compiler version" ]]
}

@test "firefox" {
  run docker container run --rm ${DOCKER_IMAGE} firefox --version
  [[ "$output" =~ "Mozilla Firefox" ]]
}

@test "lua" {
  run docker container run --rm ${DOCKER_IMAGE} lua -e 'print("シェル芸")'
  [ "$output" = "シェル芸" ]
}

@test "php" {
  run docker container run --rm ${DOCKER_IMAGE} php -r 'echo "シェル芸\n";'
  [ "$output" = "シェル芸" ]
}

@test "cureutils" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "cure girls | head -1"
  [ "$output" = "美墨なぎさ" ]
}

@test "matsuya" {
  run docker container run --rm ${DOCKER_IMAGE} matsuya
  [ $status -eq 0 ]
}

@test "takarabako" {
  run docker container run --rm ${DOCKER_IMAGE} takarabako
  [ $status -eq 0 ]
}

@test "snacknomama" {
  run docker container run --rm ${DOCKER_IMAGE} snacknomama
  [ $status -eq 0 ]
}

@test "rubipara" {
  run docker container run --rm ${DOCKER_IMAGE} rubipara kashikoma
  [ "${lines[0]}"  = '                 ／^v ＼'                                      ]
  [ "${lines[1]}"  = '               _{ / |-.(`_￣}__'                               ]
  [ "${lines[2]}"  = "        _人_  〃⌒ ﾝ'八{   ｀ノト､\`ヽ"                           ]
  [ "${lines[3]}"  = '        `Ｙ´  {l／ / /    / Ｖﾉ } ﾉ    (     Kashikoma!     )'  ]
  [ "${lines[4]}"  = '          ,-ｍ彡-ｧ Ｌﾒ､_彡ｲ } }＜く   O'                         ]
  [ "${lines[5]}"  = "         / _Uヽ⊂ﾆ{J:}  '⌒Ｖ  {  l| o"                          ]
  [ "${lines[6]}"  = "       ／  r‐='V(｢\`¨,  r=≪,/ { .ﾉﾉ"                           ]
  [ "${lines[7]}"  = '      /   /_xヘ 人 丶-  _彡ｲ ∧〉'                               ]
  [ "${lines[8]}"  = '      (  ノ¨ﾌ’  ｀^> ‐ｧｧ ＜¨ﾌｲ'                                 ]
  [ "${lines[9]}"  = "       --＝〉_丶/ﾉ { 彡' '|           Everyone loves Pripara!"  ]
  [ "${lines[10]}" = "           ^  '7^ O〉|’ ,丿"                                   ]
  [ "${lines[11]}" = '＿＿＿＿ ___ __ _{’O 乙,_r[_ __ ___ __________________________' ]
}

@test "marky_markov" {
  run docker container run --rm ${DOCKER_IMAGE} marky_markov -h
  [ "${lines[0]}" = 'Usage: marky_markov COMMAND [OPTIONS]' ]
}

@test "yq" {
  run docker container run --rm ${DOCKER_IMAGE} yq --version
  [[ "${lines[0]}" =~ "yq" ]]
}

@test "faker" {
  run docker container run --rm ${DOCKER_IMAGE} faker name
  [ $status -eq 0 ]
}

@test "sympy-python3" {
  run docker container run --rm ${DOCKER_IMAGE} python3 -c 'import sympy; print(sympy.__name__)'
  [ "$output" = "sympy" ]
}

@test "sympy" {
  run docker container run --rm ${DOCKER_IMAGE} python -c 'import sympy; print sympy.__name__'
  [ "$output" = "sympy" ]
}

@test "numpy-python3" {
  run docker container run --rm ${DOCKER_IMAGE} python3 -c 'import numpy; print(numpy.__name__)'
  [ "$output" = "numpy" ]
}

@test "numpy" {
  run docker container run --rm ${DOCKER_IMAGE} python -c 'import numpy; print numpy.__name__'
  [ "$output" = "numpy" ]
}

@test "scipy-python3" {
  run docker container run --rm ${DOCKER_IMAGE} python3 -c 'import scipy; print(scipy.__name__)'
  [ "$output" = "scipy" ]
}

@test "scipy" {
  run docker container run --rm ${DOCKER_IMAGE} python -c 'import scipy; print scipy.__name__'
  [ "$output" = "scipy" ]
}

@test "matplotlib-python3" {
  run docker container run --rm ${DOCKER_IMAGE} python3 -c 'import matplotlib; print(matplotlib.__name__)'
  [ "$output" = "matplotlib" ]
}

@test "matplotlib" {
  run docker container run --rm ${DOCKER_IMAGE} python -c 'import matplotlib; print matplotlib.__name__'
  [ "$output" = "matplotlib" ]
}

@test "xonsh" {
  run docker container run --rm ${DOCKER_IMAGE} xonsh -c 'echo シェル芸'
  [ "$output" = "シェル芸" ]
}

@test "pillow-python3" {
  run docker container run --rm ${DOCKER_IMAGE} python3 -c 'import PIL; print(PIL.__name__)'
  [ "$output" = "PIL" ]
}

@test "pillow" {
  run docker container run --rm ${DOCKER_IMAGE} python -c 'import PIL; print PIL.__name__'
  [ "$output" = "PIL" ]
}

@test "asciinema" {
  run docker container run --rm ${DOCKER_IMAGE} asciinema --version
  [[ "${lines[0]}" =~ "asciinema " ]]
}

@test "GiNZA" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | python3 -m spacy.lang.ja_ginza.cli 2>/dev/null | awk 'NR>=2{print \$3}'"
  [ "${lines[0]}" = 'シェル' ]
  [ "${lines[1]}" = '芸' ]
}

@test "egison" {
  run docker container run --rm ${DOCKER_IMAGE} egison --help
  echo "'${lines[0]}'"
  [ "${lines[0]}" = 'Usage: egison [options]' ]
}

@test "egzact" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | dupl 2"
  [ "${lines[0]}" = 'シェル芸' ]
  [ "${lines[1]}" = 'シェル芸' ]
}

@test "faker-cli" {
  run docker container run --rm ${DOCKER_IMAGE} faker-cli --help
  [ "${lines[0]}" = 'Usage: faker-cli [option]' ]
}

@test "chemi" {
  run docker container run --rm ${DOCKER_IMAGE} chemi -s H
  [ "${lines[2]}" = 'element     : Hydrogen' ]
}

@test "home-commands" {
  run docker container run --rm ${DOCKER_IMAGE} echo-sd シェル芸
  [ "${lines[0]}" = '＿人人人人人人＿' ]
  [ "${lines[1]}" = '＞　シェル芸　＜' ]
  [ "${lines[2]}" = '￣Y^Y^Y^Y^Y^Y^￣' ]
}

@test "J" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo \"'シェル芸'\" | jconsole"
  echo "'${lines[0]}'"
  [ "${lines[0]}" = '   シェル芸' ]
}

@test "trdsql" {
  run docker container run --rm ${DOCKER_IMAGE} trdsql -help
  [ "${lines[0]}" = 'Usage: trdsql [OPTIONS] [SQL(SELECT...)]' ]
}

@test "openjdk11" {
  run docker container run --rm ${DOCKER_IMAGE} javac -version
  [[ "${output}" =~ "javac " ]]
}

@test "super unko" {
  run docker container run --rm ${DOCKER_IMAGE} unko.tower 2
  [ "${lines[0]}" = '　　　　人' ]
  [ "${lines[1]}" = '　　（　　　）' ]
  [ "${lines[2]}" = '　（　　　　　）' ]
}

@test "nameko.svg" {
  run docker container run --rm ${DOCKER_IMAGE} file nameko.svg
  [ "$output" = 'nameko.svg: SVG Scalable Vector Graphics image' ]
}

@test "sayhuuzoku" {
  run docker container run --rm ${DOCKER_IMAGE} sayhuuzoku g
  [ $status -eq 0 ]
}

@test "sayhoozoku shoplist" {
  run docker container run --rm ${DOCKER_IMAGE} stat "/root/go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping/shoplist.txt"
  [ "${lines[0]}" = '  File: /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping/shoplist.txt' ]
}

@test "gron" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo '{\"s\":\"シェル芸\"}' | gron -m"
  [ "${lines[1]}" = 'json.s = "シェル芸";' ]
}

@test "pup" {
  run docker container run --rm ${DOCKER_IMAGE} pup --help
  [ "${lines[1]}" = '    pup [flags] [selectors] [optional display function]' ]
}

@test "ttyrec2gif" {
  run docker container run --rm ${DOCKER_IMAGE} ttyrec2gif -help
  [ "${lines[0]}" = 'Usage of ttyrec2gif:' ]
}

@test "owari" {
  run docker container run --rm ${DOCKER_IMAGE} owari -w 20
  [ "${lines[0]}" = '        糸冬' ]
}

@test "align" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "yes シェル芸 | head -4 | awk '{print substr(\$1,1,NR)}' | align center"
  [ "${lines[0]}" = '   シ   ' ]
  [ "${lines[1]}" = '  シェ  ' ]
  [ "${lines[2]}" = ' シェル ' ]
  [ "${lines[3]}" = 'シェル芸' ]
}

@test "taishoku" {
  run docker container run --rm ${DOCKER_IMAGE} taishoku
  [ "${lines[0]}" = '　　　代株　　　　二退こ　　　　　　' ]
}

@test "textimg" {
  run docker container run --rm ${DOCKER_IMAGE} textimg --version
  [[ "$output" =~ "textimg version " ]]
}

@test "whitespace" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo -e '   \t \t  \t\t\n\t\n     \t\t \t   \n\t\n     \t\t  \t \t\n\t\n     \t\t \t\t  \n\t\n     \t\t \t\t  \n\t\n     \t   \t\t\t\n\t\n     \t\t  \t \t\n\t\n     \t\t \t  \t\n\t\n  \n\n' | whitespace"
  [ "$output" = 'ShellGei' ]
}

@test "Open usp Tukubai" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル芸 | grep -o . | tateyoko"
  [ "$output" = 'シ ェ ル 芸' ]
}

@test "julia" {
  run docker container run --rm ${DOCKER_IMAGE} julia -e 'println("シェル芸")'
  [ "$output" = 'シェル芸' ]
}

@test "rust" {
  run docker container run --rm ${DOCKER_IMAGE} rustc --help
  [ "${lines[0]}" = 'Usage: rustc [OPTIONS] INPUT' ]
}

@test "rargs" {
  run docker container run --rm ${DOCKER_IMAGE} rargs --help
  [[ "${lines[0]}" =~ "Rargs " ]]
  [ "${lines[2]}" = 'Xargs with pattern matching' ]
}

@test "ShellGeiData" {
  run docker container run --rm ${DOCKER_IMAGE} stat /ShellGeiData/README.md
  [ "${lines[0]}" = '  File: /ShellGeiData/README.md' ]
}

@test "imgout" {
  run docker container run --rm ${DOCKER_IMAGE} imgout -h
  [ "$output" = 'usage: imgout [-f <font>]' ]
}

@test "zws" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo J+KBouKAjeKAi+KBouKAjeKAi+KAi+KAjeKAjeKBouKAjOKBouKBouKAjeKAi+KBouKAjeKAi+KAi+KAjeKAjeKAjeKAjOKBouKBouKAjeKAi+KBouKAjeKAi+KAi+KBouKAjeKAjeKAjeKBouKBouKAjeKAjeKAi+KAjeKAi+KAjeKAjeKAjeKBouKAjeKAi+KAi+KAi+KAjeKAjScK | base64 -d | ./zws -d"
  [ "$output" = 'シェル芸' ]
}

@test "osquery" {
  run docker container run --rm ${DOCKER_IMAGE} osqueryi --version
  [[ "$output" =~ 'osqueryi version ' ]]
}

@test "onefetch" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "cd /ShellGeiData && onefetch | sed $'s/\033[^m]*m//g'"
  [[ "${lines[0]}" =~ 'Project: ShellGeiData' ]]
}

@test "sushiro" {
  run docker container run --rm ${DOCKER_IMAGE} sushiro -h
  [[ "${lines[0]}" =~ 'sushiro version ' ]]
}

@test "noc" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo 部邊邊󠄓邊󠄓邉邉󠄊邊邊󠄒邊󠄓邊󠄓邉邉󠄊辺邉󠄊邊邊󠄓邊󠄓邉邉󠄎辺邉󠄎邊辺󠄀邉邉󠄈辺邉󠄍邊邊󠄓部 | mono noc -d"
  [ "$output" = 'シェル芸' ]
}

@test "bat" {
  run docker container run --rm ${DOCKER_IMAGE} bat --version
  [[ "$output" =~ "bat " ]]
}

@test "echo-meme" {
  run docker container run --rm ${DOCKER_IMAGE} echo-meme シェル芸
  [[ "$output" =~ "シェル芸" ]]
}

@test "bash 5.0" {
  run docker container run --rm ${DOCKER_IMAGE} bash --version
  [[ "${lines[0]}" =~ "GNU bash, バージョン 5.0" ]]
}

@test "awk 5.0" {
  run docker container run --rm ${DOCKER_IMAGE} /usr/local/bin/awk --version
  [[ "${lines[0]}" =~ "GNU Awk 5.0" ]]
}

@test "reiwa" {
  run docker container run --rm ${DOCKER_IMAGE} date -d '2019-05-01' '+%Ec'
  [ "$output" = '令和元年05月01日 00時00分00秒' ]
}

@test "NormalizationTest.txt" {
  run docker container run --rm ${DOCKER_IMAGE} stat NormalizationTest.txt
  [ "${lines[0]}" = '  File: NormalizationTest.txt' ]
}

@test "NamesList.txt" {
  run docker container run --rm ${DOCKER_IMAGE} stat NamesList.txt
  [ "${lines[0]}" = '  File: NamesList.txt' ]
}

@test "ke2dair" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "echo シェル 芸 | ke2daira"
  [ "$output" = 'ゲェル シイ' ]
}

@test "man" {
  run docker container run --rm ${DOCKER_IMAGE} bash -c "man シェル芸 |& cat"
  [ "$output" = 'シェル芸 というマニュアルはありません' ]
}
