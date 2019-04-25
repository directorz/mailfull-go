# Mailfull の使い方

  * ユーザは `mailfull` にて実行してください。
  * すべてのコマンドは、`-h` オプションで簡単な使用方法を出力して終了します。
  * すべてのコマンドは、`-v` オプションでバージョンを出力して終了します。
  * `-n` オプションでデータベースをアップデートしません。(直ちに変更を行いません)  


## ドメイン

### ドメインの追加

    $ mailfull2 domainadd example.com

  `example.com` が追加され、`postmaster@example.com` が追加されます。 
  `postmaster@example.com` のパスワードは付与されませんので `userpasswd` コマンドで対応してください。 

### ドメインの削除

    $ mailfull2 domaindel example.com
  `example.com` が削除されます。 
  ドメインディレクトリは同階層にドット付きでバックアップされます。 

### ドメインのリストアップ

    $ mailfull2 domains

  設定されているドメインがリストアップされます。 


## ユーザ

### ユーザの追加

    $ mailfull2 useradd user@example.com

  `user@example.com` が追加されます。 
  パスワードは付与されませんので `userpasswd` コマンドで対応してください。 

  SMTP-AUTH, POP/IMAP のユーザ名は、`user@example.com` となります。

### ユーザの削除

    $ mailfull2 userdel user@example.com

  `postmaster` は削除できません。 
  ユーザディレクトリは同階層にドット付きバックアップされます。

### パスワードの変更

    $ mailfull2 userpasswd user@example.com
    Enter new password for user@example.com:
    Retype new password:

  `user@example.com` のパスワードを設定、変更します。

### ユーザのリストアップ

    $ mailfull2 users example.com

  `example.com` のユーザがリストアップされます。

### パスワードのチェック

    $ mailfull2 usercheckpw user@example.com
    Enter password for user@example.com: 
    The password you entered is correct.
    $ echo $?
    0

    $ mailfull2 usercheckpw user@example.com
    Enter password for user@example.com: 
    The password you entered is incorrect.
    $ echo $?
    1

  パスワードは引数で与えることもできます。 
  パスワードの正誤に応じたメッセージが表示されます。 
  コマンドの戻り値として、正しい場合 0、違う場合 1 が戻ります。 


## エイリアス

### エイリアスの新設

    $ maulfull2 aliasuseradd aliasname@example.com dest@example.org[,...]

  `aliasname@example.com` 宛のメールを `dest@example.org` へ転送するエイリアスを作成します。

### エイリアスの編集

    $ maulfull2 aliasusermod aliasname@example.com dest2@example.org[,...]

  `aliasname@example.com` 宛のメールを `dest2@example.org` へ転送するエイリアスへ書き換えます。

### エイリアスの解除

    $ mailfull2 aliasuserdel aliasname@example.com 

  `aliasname@example.com` のエイリアス設定を削除します。 

### エイリアスのリストアップ

    $ mailfull2 aliasusers example.com

  ドメインに設定されているエイリアスのリストが出力されます。


## メーリングリスト

### メーリングリストの一覧

    $ mailfull2 mailinglists example.com

  `example.com` のMLがリストアップされます。

### メーリングリストの新規作成

    $ mailfull2 mailinglistadd example.com ml user@example.com

  `ml@example.com` というメーリングリストを作成し、メンバーに `user@example.com` を含めます。

### 既存のメーリングリストに配信先を追加

    $ mailfull2 mailinglistuseradd example.com ml user2@example.com

  `ml@example.com` というメーリングリストに、メンバーに `user2@example.com` を追加します。

### 既存のメーリングリストの配信先を確認

    $ mailfull2 mailinglistusers example.com ml

  `ml@example.com` というメーリングリストのメンバーを確認します。

### 既存のメーリングリストから配信先を削除

    $ mailfull2 mailinglistuserdel example.com ml user@example.com

  `ml@example.com` というメーリングリストから、メンバー `user@example.com` を削除します。

### メーリングリストの削除

    $ mailfull2 mailinglistdel example.com ml

  `ml@example.com` というメーリングリストを削除します。

## キャッチオール

### キャッチオールの設定

    $ mailfull2 catchallset example.com user

  `example.com` の全てのユーザ宛のメールを、 
  ユーザが存在しなければ、`user@example.com` が受け取るようになります。 

### キャッチオールの解除

    $ mailfull2 catchalldel example.com

  `example.com` に設定されているキャッチオールを解除します。

### キャッチオールの取得

    $ mailfull2 catchall example.com

  `example.com` にキャッチオールが設定されていれば、出力します。


## エイリアスドメイン

### エイリアスドメインの追加

    $ mailfull2 aliasdomainadd alias.example.com example.com

  `example.com` 宛のメールが、`alias.example.com` でも  
  受け取れるようになります。

### エイリアスドメインの解除

    $ mailfull2 aliasdomaindel alias.example.com

  エイリアスドメインを解除します。

### エイリアスドメインのリストアップ

    $ mailfull2 aliasdomains example.com

  `example.com` に設定されているエイリアスドメインをリストアップします。


## その他

### commit

    $ ./commit

  `/home/mailfull/domains` 以下から設定を生成し、 
  `/home/mailfull/etc` 以下の設定ファイルにまとめ、各種データベースを作成します。

   
