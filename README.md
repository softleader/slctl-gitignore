# slctl-gitignore

The [slctl](https://github.com/softleader/slctl) plugin to fetches a .gitignore from [gitignore.io](https://gitignore.io/)

> 建立或增加更多的 .gitignore

## Install

```sh
$ slctl plugin install github.com/softleader/slctl-gitignore
```

## Usage

從 [gitignore.io](https://gitignore.io/) 取得指定類別 .gitignore 的內容

```sh
$ gitignore go intellij+all
```

若不指定 ignore 類別, 會有一個互動式選單讓你更方便的做選擇

```sh
$ gitignore
```

使用 `--outout` 指定路徑, 可以將內容儲存成檔案而不是印在 console
如果指定的路徑檔案已經存在, 將會將內容接續在原本的檔案中

```sh
$ gitignore -o /path/to/.gitignore
```

如果指定的路徑是目錄, 將會以 `.gitignore` 為預設檔名儲存檔案

```sh
$ gitignore -o .
```
