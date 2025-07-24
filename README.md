# 2025_School_Experiment
## 命名規則
なるべくこれに沿う。  
まあ柔軟にいきましょう。短めが望ましい。

[【Go言語】Goらしい命名規則](https://qiita.com/rapirapi/items/4c8b41996aa3071a5aeb)

## コミットメッセージ
「ニ字熟語 : メッセージ」の形。  
例：「例文 : これはコミットメッセージの例文です。」

## コメントアウト
関数、構造体にはコメントアウトをつける。それ以外のものは処理で一見すると難しそうとかであれば書く。  
書き過ぎには注意、可読性下がるので。

## ディレクトリ構造
SHISO_CHECKER/
├── README.md              # プロジェクトの概要
├── main.go                # エントリーポイント
├── go.mod                 # Goモジュール設定
├── /db/                   # DB関連群
│   └── db.go              # DB処理
├── /handlers/             # ルーティングされたハンドラ群
│   ├── auth.go            # 認証関連
│   ├── middleware.go      # JWT認証用のGinミドルウェア
│   └── users.go           # ユーザーのHTTPリクエストまとめ
├── /models/               # データ構造
│   └── users.go           # ユーザーのデータ形式
└── /utils/                # 汎用処理群
    ├── age.go             # 年齢計算周り
    └── jwt.go             # JWTの汎用処理