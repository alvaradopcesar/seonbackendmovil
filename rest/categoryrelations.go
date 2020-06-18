package rest

import (
	"log"
	"strconv"

	"seonbackendmovil/utils"

	"github.com/gin-gonic/gin"
)

// Strcategoryrelations sdatos de retorno
type Strcategoryrelations struct {
	ID              int64  `db:"id" json:"id"`
	CategoryID      int64  `db:"categoryId" json:"categoryId"`
	Name            string `db:"name" json:"name"`
	ChildCategoryID int64  `db:"childCategoryId" json:"childCategoryId"`
	NameChild       string `db:"nameChild" json:"nameChild"`
}

// GetCategoryrelationsChildsAndParents trae padres e hijos de CategoryRelations
func GetCategoryrelationsChildsAndParents(c *gin.Context) {
	log.Println("************** START [ /category/parents (GetCategoryrelationsChildsAndParents) ] **************")
	var strcategoryrelations []Strcategoryrelations
	_, err := utils.DbmapMySQLInkafarma.Select(&strcategoryrelations,
		`
		select cr.id , cr.categoryId,c1.name , cr.childCategoryId , c2.name as nameChild
  		  from categoryrelations cr
	   		   inner join category c1 on cr.categoryId = c1.id
       		   inner join category c2 on cr.childCategoryId = c2.id
 	  order by cr.categoryId, cr.childCategoryId
	`)

	log.Println("************** STOP [ /getcategorybyparents/{id} (GetCategoryrelationsChildsAndParents) ] **************")
	if err == nil {
		c.JSON(200, strcategoryrelations)
	} else {
		c.JSON(404, gin.H{"error": err.Error()})
	}
}

// GetCategoryRelationsByID devuelve categorias asociadas
func GetCategoryRelationsByID(c *gin.Context) {
	log.Println("************** START [ /getcategorybyparents/{id} (GetCategoryRelationsByID) ] **************")
	sid := c.Params.ByName("id")
	id, _ := strconv.ParseInt(sid, 10, 64)

	var strcategoryrelations []Strcategoryrelations
	log.Println(id)
	data := getCategoryParent(id)
	log.Println("================")
	for ind, dataRow := range data {
		log.Printf("%d: %v\n", ind, dataRow)
		strcategoryrelationsRow, err := getCategoryrelationsByID(dataRow.ID)
		if err == nil {
			strcategoryrelations = append(strcategoryrelations, strcategoryrelationsRow)
		} else {
			c.JSON(404, gin.H{"error": err.Error()})
		}
	}
	log.Println("================")
	data = getCategoryChild(id)
	for ind, dataRow := range data {
		log.Printf("%d: %v\n", ind, dataRow)
		strcategoryrelationsRow, err := getCategoryrelationsByID(dataRow.ID)
		if err == nil {
			strcategoryrelations = append(strcategoryrelations, strcategoryrelationsRow)
		}
	}
	log.Println("************** STOP [ /category/parents/:id (GetCategoryRelationsByID) ] **************")
	c.JSON(200, strcategoryrelations)
}

// getCategoryrelationsByID df
func getCategoryrelationsByID(id int64) (Strcategoryrelations, error) {

	var strcategoryrelations Strcategoryrelations
	err := utils.DbmapMySQLInkafarma.SelectOne(&strcategoryrelations,
		`select cr.id , cr.categoryId,c1.name , cr.childCategoryId , c2.name as nameChild
  		  from categoryrelations cr
	   		      inner join category c1 on cr.categoryId = c1.id
				  inner join category c2 on cr.childCategoryId = c2.id
		 where cr.id = ?		  
	   order by cr.categoryId, cr.childCategoryId`, id)
	if err != nil {
		return strcategoryrelations, err
	}
	return strcategoryrelations, nil
}

// StrGetCategotyParameter w
type StrGetCategotyParameter struct {
	ID              int64 `db:"id" json:"id"`
	CategoryID      int64 `db:"CategoryID" json:"CategoryID"`
	ChildCategoryID int64 `db:"ChildCategoryID" json:"ChildCategoryID"`
}

func getCategoryParent(id int64) []StrGetCategotyParameter {
	var strGetCategotyParameters []StrGetCategotyParameter
	_, err := utils.DbmapMySQLInkafarma.Select(&strGetCategotyParameters,
		` SELECT id,CategoryID,ChildCategoryID FROM categoryrelations 
		   where CategoryID = ?`,
		id)
	if err != nil {
		return nil
	}
	var row []StrGetCategotyParameter
	for ind, dataRow := range strGetCategotyParameters {
		row = append(row, dataRow)
		log.Printf("%d: %v\n", ind, dataRow)
		strGetCategotyParameters2 := getCategoryParent(dataRow.ChildCategoryID)
		for _, dataRow2 := range strGetCategotyParameters2 {
			row = append(row, dataRow2)
		}
	}
	return row
}

func getCategoryChild(id int64) []StrGetCategotyParameter {
	var strGetCategotyParameters []StrGetCategotyParameter
	_, err := utils.DbmapMySQLInkafarma.Select(&strGetCategotyParameters,
		` SELECT id,CategoryID,ChildCategoryID FROM categoryrelations
		   where ChildCategoryID = ?`,
		id)
	if err != nil {
		return nil
	}
	var row []StrGetCategotyParameter
	for ind, dataRow := range strGetCategotyParameters {
		row = append(row, dataRow)
		log.Printf("%d: %v\n", ind, dataRow)
		strGetCategotyParameter2 := getCategoryChild(dataRow.CategoryID)
		for _, dataRow2 := range strGetCategotyParameter2 {
			row = append(row, dataRow2)
		}
	}
	return row
}
