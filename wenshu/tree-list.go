package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

func treeListDemo() {
	vm := otto.New()
	x, _ := vmRunS(vm, treelist)
	x, _ = vmRunS(vm, fmt.Sprintf(`JSON.stringify(%s)`, x))
	fmt.Println(x)
}

const treelist = `"[{ \u0027title\u0027:\u0027关键词\u0027,\u0027key\u0027:\u0027keyword\u0027,\u0027path\u0027:\u0027案/key\u0027},{ \u0027title\u0027:\u0027案由\u0027,\u0027key\u0027:\u0027caseType\u0027,\u0027path\u0027:\u0027案/caseType\u0027},{ \u0027title\u0027:\u0027法院\u0027,\u0027key\u0027:\u0027court\u0027,\u0027path\u0027:\u0027案/court\u0027},{ \u0027title\u0027:\u0027裁判年份\u0027,\u0027key\u0027:\u0027trialYear\u0027,\u0027path\u0027:\u0027案/trialYear\u0027},{ \u0027title\u0027:\u0027审理程序\u0027,\u0027key\u0027:\u0027trialRound\u0027,\u0027path\u0027:\u0027案/trialRound\u0027},{ \u0027title\u0027:\u0027文书性质\u0027,\u0027key\u0027:\u0027judgeType\u0027,\u0027path\u0027:\u0027案/judgeType\u0027}]"`
