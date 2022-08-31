package cmd

import "fmt"

func Execute() {
	fmt.Println("# 認証を行う")
	fmt.Println("## リポジトリの一覧のループ")
	{
		fmt.Println("### 日付の一覧のループ")
		{
			fmt.Println("#### prの一覧を取得")
			fmt.Println("#### approveを取得")
			fmt.Println("#### commentを取得")
			fmt.Println("#### prの数を取得")
		}
	}
}
