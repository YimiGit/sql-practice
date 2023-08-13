package grpcService

import (
	"common/config"
	model "common/model/practiceModel"
	"common/proto/practiceProto"
	"context"
	"log"
	"strings"
	"sync"
)

// PracticeService 答题相关rpc服务
type PracticeService struct {
	practiceProto.UnimplementedPaperServiceServer
}

// LevelList 获取难度类型列表
func (*PracticeService) LevelList(c context.Context, req *practiceProto.LevelListRequest) (*practiceProto.LevelListResponse, error) {
	var levelList []*practiceProto.Level
	if result := config.DB.Model(&model.PracticeLevel{}).Order("type asc").Find(&levelList); result.Error != nil {
		return nil, result.Error
	}
	return &practiceProto.LevelListResponse{LevelList: levelList}, nil
}

// PaperList 获取所有试卷列表
func (*PracticeService) PaperList(c context.Context, req *practiceProto.PaperListRequest) (*practiceProto.PaperListResponse, error) {
	var paperList *[]*practiceProto.Paper
	if result := config.DB.Model(&model.PracticePaper{}).Find(&paperList); result.Error != nil {
		return nil, result.Error
	}
	paperMap := make(map[int32]*practiceProto.PaperList)
	for _, v := range *paperList {
		if _, exist := paperMap[v.Type]; !exist {
			paperMap[v.Type] = &practiceProto.PaperList{}
		}
		paperMap[v.Type].PaperList = append(paperMap[v.Type].PaperList, v)
	}
	return &practiceProto.PaperListResponse{PaperMap: paperMap}, nil
}

// QuestionList 通过PaperID获取试卷的所有题目
func (*PracticeService) QuestionList(c context.Context, req *practiceProto.QuestionListRequest) (*practiceProto.QuestionListResponse, error) {
	var questionList []*practiceProto.Question
	if result := config.DB.Model(&model.Practice{}).Where("paper_id = ?", req.PaperID).Find(&questionList); result.Error != nil {
		return nil, result.Error
	}
	return &practiceProto.QuestionListResponse{QuestionList: questionList}, nil
}

// TableStruct 通过PaperID获取 该试卷所需 表结构
func (*PracticeService) TableStruct(c context.Context, req *practiceProto.TableStructRequest) (*practiceProto.TableStructResponse, error) {
	var structString string
	if result :=
		config.DB.Model(&model.PracticePaper{}).
			Select("table_struct").
			Where("id = ?", req.PaperID).
			First(&structString); result.Error != nil {
		return nil, result.Error
	}

	columnsMap := make(map[string]*practiceProto.TableStruct)

	//该试卷需要用到的表
	structs := strings.Split(structString, ",")

	//协调携程
	group := sync.WaitGroup{}
	group.Add(len(structs))

	//并发查询表结构
	for _, s := range structs {
		go func(ss string) {
			defer func() {
				group.Done()
			}()
			//查询表结构
			columns, err := config.DB.Migrator().ColumnTypes(ss)
			if err != nil {
				log.Println("查询表结构失败")
				return
			}
			//组装数据结构
			var columnComments practiceProto.TableStruct
			var columnCommentList []*practiceProto.ColumnComment
			columnComments.ColumnCommentList = columnCommentList
			for _, column := range columns {
				value, _ := column.Comment()
				columnComments.ColumnCommentList = append(columnComments.ColumnCommentList, &practiceProto.ColumnComment{Field: column.Name(), Comment: value})
			}
			columnsMap[ss] = &columnComments
		}(s)
	}
	group.Wait()
	return &practiceProto.TableStructResponse{TableStructMap: columnsMap}, nil
}
