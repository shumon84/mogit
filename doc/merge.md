# merge の仕組み
> 参照 https://www.atlassian.com/ja/git/tutorials/using-branches/merge-strategy

実は git merge はいつも同じアルゴリズムを使っているわけではない。
指定された strategy によって使い分けている(指定されなかった場合は、最適な strategy を自動で使ってくれる))

- recursive
  - 一番よく使うやつ(というかほとんどの人はこれしか使ったことがない)
- resolve
  - 3-way merge した 2 つの HEAD の解決のみに使う
  - 十字マージを検出できるて、高速かつ安全とされている
- octopus
  - 3-way merge 時にデフォルトで使われる strategy
  - 
- ours