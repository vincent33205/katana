package utils

import "testing"

// TestTransformIndex 測試 TransformIndex 函數的各種邊界情況和正常情況
// 涵蓋以下測試場景：
//   - 空切片處理
//   - 第一個元素訪問
//   - 範圍內索引轉換
//   - 下界鉗位（索引過小）
//   - 上界鉗位（索引過大）
func TestTransformIndex(t *testing.T) {
	// 測試空切片的情況
	t.Run("empty slice", func(t *testing.T) {
		if got := TransformIndex([]int{}, 5); got != 0 {
			t.Fatalf("空切片應返回 0，實際得到 %d", got)
		}
	})

	// 定義測試用例結構
	type testCase struct {
		name  string   // 測試用例名稱
		arr   []string // 測試陣列
		index int      // 輸入索引
		want  int      // 期望結果
	}

	// 測試用例集合
	cases := []testCase{
		{name: "first element", arr: []string{"a", "b", "c"}, index: 1, want: 0}, // 第一個元素
		{name: "in range", arr: []string{"a", "b", "c"}, index: 2, want: 1},      // 範圍內索引
		{name: "clamp low", arr: []string{"a", "b", "c"}, index: 0, want: 0},     // 下界鉗位
		{name: "clamp high", arr: []string{"a", "b", "c"}, index: 5, want: 2},    // 上界鉗位
	}

	// 執行所有測試用例
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := TransformIndex(tc.arr, tc.index); got != tc.want {
				t.Fatalf("期望 %d，實際得到 %d", tc.want, got)
			}
		})
	}
}
