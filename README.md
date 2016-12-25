# ghb

ghb は Github リポジトリを操作する CLI です。

## Description

* 新しいリポジトリを作成できます。
* 新しい Issue を立てることができます。

## Motivation

* そろそろ Go ツールを作れるようになってきた。
* ブラウザからリポジトリ作ったり Issue 登録するのだるい。端末上で済ませたい。

## Installation

[Releases ページ](https://github.com/yuta-masano/ghb/releases)からダウンロードしてください。

## Usage

### `git add ...`

#### 新規リポジトリ作成

```
$ ghb add repository test-repo -d 'This is a test repository.' -u 'https://example.com'
```

* 引数として作成したいリポジトリ名を指定します。
* フラグオプションでリポジトリの short description と URL を追加できます。
* リポジトリの作成に成功すると、標準出力にリポジトリ URL を表示します。

#### 新規 Issue 作成

```
$ ghb add issue test-repo -l bug
```

* 引数として Issue を作成したいリポジトリ名を指定します。
* フラグオプションで Issue に
* リポジトリの作成に成功すると、標準出力に Issue 番号を表示します。

### Option

## License

The MIT License (MIT)

## Thanks to

## Author

[Yuta MASANO](https://github.com/yuta-masano)

## Development

### セットアップ

```
$ # 1. リポジトリを取得。
$ go get -v -u -d github.com/yuta-masano/ghb

$ # 2. リポジトリディレクトリに移動。
$ cd $GOPATH/src/github.com/yuta-masano/ghb

$ # 3. 開発ツールと vendor パッケージを取得。
$ make deps-install

$ # 4. その他のターゲットは help をどうぞ。
$ make help
USAGE: make [target]

TARGETS:
help           show help
...
```

### リリースフロー

see: [yuta-masano/dp#リリースフロー](https://github.com/yuta-masano/dp#%E3%83%AA%E3%83%AA%E3%83%BC%E3%82%B9%E3%83%95%E3%83%AD%E3%83%BC)
