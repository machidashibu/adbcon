# [ADR-0003] [implement-multi-handler-for-slog-by-DIY]
# [ADR-0003] [slogのマルチハンドラを自前でコーディングする]
<!--
基本ルール：
* 英語のタイトルを日本語のタイトルを併記する
* 英語のタイトルはそのままファイル名に使えるようにハイフンで連結する。
* 見出しは英語を使うが、本文は日本語で書く。（変に英語で書いて自分で読み返しにくくしないように）
* 個人プロジェクトなのでStatusは基本Acceptedで始まる。（Proposedは基本使われない。Rejectedは考えたけどやめたこととして使う場合がある）
-->


* **Status:** Accepted
* **Date:** 2026-9-24
<!--
* **ステータス:** [提案中 / 承認済 / 拒否 / 非推奨 / 撤回（ADR-XXXXに置換）]
-->

## Context
<!-- 
## 1. 背景と課題 (Context)
* なぜこの決定が必要なのか？
* 解決したい問題や制約条件（納期、コスト、技術的制約など）は何か？
-->
* 出力先ごとにデバッグレベルを分けたい。
  * 標準出力には常に`Info`以上のみ出力
  * ログファイルには常に`Debug`以上を出力
* 上記で、開発中に`Info`レベルのログ出力の妥当性を目視で確認できる。
  * verboseモードなどで`Debug`も標準出力に出していると`Info`の出力が妥当かを見逃しがち。
* かつ問題解析時にログファイルを見ることで`Debug`も出力されているのでスムーズに行える。

## Options Considered
1. **(Chosen)** 複数の`slog.Handler`を束ねて処理するMultiHandlerを自前で実装する。
2. `samber/slog-multi`を使用する。
  * → See: https://github.com/samber/slog-multi
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
* 1を選択する。
### Rationale
* 今回は勉強もかねてslog.Handlerを自前で実装する。

## Consequences
<!--
* **メリット（プラスの影響）:** 
  * 何が良くなるか？
* **デメリット・リスク（マイナスの影響）:** 
  * どのようなトレードオフや課題が残るか？
-->
### Positive Impact (Pros):
* What benefits or advantages do we gain?
### Negative Impact / Trade-offs (Cons):
* What are the drawbacks, risks, or technical debt we accept?

