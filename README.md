# log-analyzer

"log-analyzer" can analysis  apache log and output the number of access per time and client IP.

apache ログを分析し、時間帯毎のアクセス数とクライアントIP毎のアクセス数を出力するプログラム。

## Install
実行環境として"Docker"のコンテナ環境を採用した。
* 実行環境
  - ホストOS : OSX
  - Docker version 17.03.1-ce
  - コンテナイメージ
    - ubuntu:16.04
    - go version go1.6.2 linux/amd64


1. Install docker  
  * [docker docs](https://docs.docker.com/engine/installation/)
  を参照する
2. Build container image  
[image name]はビルドするコンテナイメージの名前になるので、任意に指定する。
   >  
   $ git clone https://github.com/makima333/log-analyzer  
   $ cd log-analyzer    
   $ docker build -t [image name] .

3. Run container
  * case : analysis log  
  ログ解析する時、ログがあるディレクトリをコンテナのディレクトリにマウントします。
  > $ docker run -it --rm -v [log directory PATH]:/root/log [image name]

  * case : development  
  log-analyzerを改良する際には、~/log-analyzer/goをコンテナのディレクトリにマウントします。
  > $ docker run -it --rm -v ~/log-analyzer/go:/root/go [image name]

  コンテナ内のコンソールが表示されます。(*65c8dfdcf76f* はコンテナIDなので実行毎で変化します)  
  ``root@65c8dfdcf76f:~/log#  ``

4. Exec log-analyzer

```
root@65c8dfdcf76f:~/log# log-analyzer
NAME:
   parse apache log - A new cli application

USAGE:
   log-analyzer [global options] command [command options] [arguments...]

VERSION:
   1.0

COMMANDS:
     time, t  echo time_map
     host, h  echo host_rank
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --start value, -s value  
   --end value, -e value    
   --help, -h               show help
   --version, -v            print the version  
```

## Usage
 * SubCommand
  - time:時間毎のアクセス数を表示します。
    > # log-analyzer time [file name]
  - host:クライアントIP毎のアクセス数を表示します。
    > # log-analyzer host [file name]

 * Option
  - "-s": ログ解析開始日の指定　例:2016-01-01
    > # log-analyzer -s 2016-01-01 time [file name]

  - "-e": ログ解析修了日の指定 "-s" と併用も可能
    > # log-analyzer -e 2017-01-01 time [file name]


 * Multi log file
 複数ファイルを指定することも可
 > # log-analyzer time [file name01] [file name02]

 * exapmle
 "test01.log"と"test02.log"の2016年4月1日から2017年2月1日までのクライアントIP毎のアクセス数が知りたい場合
 > $ log-analyzer -s 2016-04-01 -e 2017-02-01 host test01.log test02.log

## Other
  * flow chart
  ①ファイル読み込み　②1行取得　③1行分のlogの各情報を構造体(line)に格納しバッファ(lfactor)に書き込む　④②に戻る  
  これと並列して以下の作業を行う
  ①バッファの読み込み ②日付の判定 ③集計用の変数に格納

  * faster (high spec cpu)
  マルチコアのCPUで実行する場合はファイルを分割して実行した方が早い。
    - 1万行のログファイルを分析した場合の実行時間　0m3.160s
    - 5千行のログファイルを2つ指定した場合の実行時間　0m2.236s
