# [ADR-0002] [use-middleware-of-opai-codegen-to-validate-API]
# [ADR-0002] [APIのバリデーションはopai-codegenのミドルウェアを使用する]
<!--
基本ルール：
* 英語のタイトルを日本語のタイトルを併記する
* 英語のタイトルはそのままファイル名に使えるようにハイフンで連結する。
* 見出しは英語を使うが、本文は日本語で書く。（変に英語で書いて自分で読み返しにくくしないように）
* 個人プロジェクトなのでStatusは基本Acceptedで始まる。（Proposedは基本使われない。Rejectedは考えたけどやめたこととして使う場合がある）
-->


* **Status:** Accepted
* **Date:** 2026-9-23
<!--
* **ステータス:** [提案中 / 承認済 / 拒否 / 非推奨 / 撤回（ADR-XXXXに置換）]
-->

## Context
<!-- 
## 1. 背景と課題 (Context)
* なぜこの決定が必要なのか？
* 解決したい問題や制約条件（納期、コスト、技術的制約など）は何か？
-->
* 当初はoapi-codebenの生成コードで使用する型（`XxeParams`や`XxxRequestBody`など）に`Validate`を実装してハンドラーの先頭で呼び出していた。
* しかし、`type XxxRequestBody = []XxxList`のように出力される場合`validate`を実装できない。
* あまりOpen APIの定義に変な記述を追加したり、構造を歪めずにバリデーションを行うように変更する。

## Options Considered
1. **(Chosen)** oapi-codegenプロジェクトが提供してるミドルウェアを使用する。
  * refs: [oapi-codegen/echo-middleware](https://github.com/oapi-codegen/echo-middleware)
2. `internal/adapter/model`に`Validate`実装用の型を用意して`x-go-type`でopenapiの定義と関連付ける。
3. `requestBody`から`$ref`する場合はすべて`type: object`にする。
  * ↑ で必ずメソッドが追加可能な`type`でコードが生成される。
<!--
## 2. 検討した選択肢 (Options Considered)
* **選択肢A（採用案）:** [簡単な説明]
* **選択肢B:** [簡単な説明]
* **選択肢C:** [簡単な説明]
-->

## Decision
<!--
* どの選択肢を採用するか？
* **決定理由:** なぜその選択肢を選んだのか？
-->
* 1番を選択する。
### Rationale
* どちらかというと消去法的な結論（だが納得しやすいので）
* 2番は以下の点が気になる。
  * `x-go-type`という異物（APIを単に読みたいだけの人に不要な）がOpen APIの定義に入り込む。
  * `x-go-type`で指定するための型は`internal/adapter/model` （domainの実装）に実装することになるが`import`がループ参照になる可能性がある。（oapi-codegenの生成コードがアプリ内の`internal/adapter/model`を参照することになるので）
* 3番目はOpen APIの定義をゆがめるので論外。（読む人に「なぜこんなことを？という疑問を与える」）
* 生成コードのパッケージから自前コードを除去できる。（ちょっと気持ち悪かった）

## Consequences
<!--
* **メリット（プラスの影響）:** 
  * 何が良くなるか？
* **デメリット・リスク（マイナスの影響）:** 
  * どのようなトレードオフや課題が残るか？
-->
### Positive Impact (Pros):
* goのコードからバリデーションのコードをなくせる。
* 正しくエラーにするためにより正確に、明確にOpen APIの定義がかかれる。
### Negative Impact / Trade-offs (Cons):
* Open APIの定義が不十分だと不正なデータがすり抜けてきてしまう。
* そのためのテストを書くことが少し手間
  * 実際に動作しているサーバに結合テストとしてテスト用クライアントまたはcurlでエラーデータを流し込む。
  * `internal/infra/server`の単体テストでもできるが、前述の結合テストっぽい方法のほうが楽かも。
* ハンドラーの先頭に`Validate`がないのはちょっと気持ち悪い
  * 今までその形式に慣れているので。
  * あと、実際の`Validate`の処理が`middleware`に隠蔽されているので
    * この点は正しいOpen APIの定義と、正しいテストの作成を促すという意味では単なるリスクではない。
