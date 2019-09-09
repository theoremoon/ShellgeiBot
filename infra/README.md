# infra

ShellgeiBotをdaemon化するための環境構築手順と動作確認方法を記載。

## 目的

ShellgeiBotをdaemon化し、異常終了した際に自動でプロセスを復旧する仕組みを導入する。

## 依存ツール

ShellgeiBotプロセスの自動起動に以下のツールを使用します。

* supervisor

その他、アプリをビルドするために以下のコマンドを使用します。

* git
* make

## 環境構築

Ubuntu環境での環境構築手順は以下のとおりです。

```bash
sudo apt update -y
sudo apt install -y supervisor
sudo git clone https://github.com/theoremoon/ShellgeiBot /opt/ShellgeiBot

# ShellgeiBotアプリのビルド
cd /opt/ShellgeiBot && make build

# 設定ファイルのリンク
sudo ln -sfn /opt/ShellgeiBot/infra/etc/supervisor/conf.d/ShellgeiBot.conf /etc/supervisor/conf.d/ShellgeiBot.conf
# ログ出力先の作成
sudo mkdir /var/log/ShellgeiBot

# supervisor のステータス確認
sudo service supervisor status
## -> running でなければ以下を実行
sudo service supervisor start
## -> running だった場合は再起動
sudo service supervisor restart

# supervisor のステータス確認
sudo service supervisor status
## -> running になっていることを確認

# daemonとしてShellgeiBotが存在することを確認
ps aux | grep ShellgeiBot
```

## ShellgeiBot daemonの自動復旧の確認

supervisordがShellgeiBotプロセスの常駐を監視しているため、プロセスが死んだ場合に自動起動します。
以下の手順でsupervisorによってプロセスが自動で復旧することを確認できます。

```bash
# ShellgeiBotプロセスIDを確認
ps aux | grep ShellgeiBot

# 前述の方法で確認したプロセスIDを指定
kill 'プロセスID'

# ShellgeiBotプロセスが存在しており、かつプロセスIDが変化していることを確認
ps aux | grep ShellgeiBot
```

## ログファイル

何らかの問題が発生したことでログを調査する場合は、以下の箇所を確認してください。

| パス | 説明 |
|------|------|
| /var/log/supervisor/supervisord.log | ShellgeiBot daemonを管理するサービス自体のログ |
| /var/log/ShellgeiBot/ShellgeiBot.log | ShellgeiBotの標準出力の記録されるログファイル |
| /var/log/ShellgeiBot/ShellgeiBot_error.log | ShellgeiBotの標準エラー出力の記録されるログファイル |

## daemon化前の動作確認

いきなりコマンドをホストPC上で実行する前にリハーサルしたい場合は以下のコマンドを実行する。

```bash
make build

# コンテナに入る
make start

# ダミーのスクリプトをShellgeiBotコマンドの代わりに配置
cp infra/ShellgeiBot /opt/ShellgeiBot/

# supervisorが起動していないことを確認
service supervisor status
# 起動
service supervisor start
# supervisorが起動していることを確認
service supervisor status

# daemon化したShellgeiBotのプロセスを確認
ps aux | grep ShellgeiBot

# ダミーのコマンドはログをここに出力するだけ
tail -f /var/log/ShellgeiBot/ShellgeiBot.log

# 試しにプロセスを落とす
kill 'プロセスID'

# プロセスIDが変わっていることを確認
ps aux | grep ShellgeiBot
```

<!-- vim: set tw=0 nowrap: -->
