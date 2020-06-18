package rest

import (
	"log"
	"seonbackendmovil/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRelProducts nuevo
func GetRelProducts(c *gin.Context) {
	log.Println("************** START [ /getrelproducts (GetRelProducts) ] **************")

	type Params struct {
		ID string `form:"id" json:"id"`
		// getrelproducts?id=001807
	}
	var params Params
	if c.Bind(&params) != nil {
		log.Println("Usted debera ingresar todos los parametros ")
		c.JSON(500, gin.H{"error": "Usted debera ingresar todos los parametros "})
		return
	}

	type ProductRelated struct {
		ProductRelatedTypeid int64  `db:"productrelatedtypeid" json:"ProductRelatedTypeid"`
		ProductRelated       string `db:"productRelated" json:"productRelated"`
	}
	query :=
		`select pr.ProductRelatedTypeid ,
			    pr.productRelated
			from ProductRelated pr 
		  where pr.productId = ?`
	var productRelateds []ProductRelated
	_, err := utils.DbmapMySQLInkafarma.Select(&productRelateds, query, params.ID)
	if err != nil {
		log.Println(query)
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	type CategoryList struct {
		ID   int64  `db:"id" json:"id"`
		Name string `db:"Name" json:"name"`
	}

	type ProductActiveComponent struct {
		Name          string `db:"name" json:"name"`
		Concentration string `db:"concentration" json:"concentration"`
	}

	type ImageLists struct {
		URL string `db:"imagePath" json:"url"`
	}

	type ReturnQuery struct {
		// Keyword string `db:"keyword" json:"keyword"`
		ID                 string         `db:"id" json:"id"`
		Name               string         `db:"Name" json:"name"`
		LongDescription    string         `db:"LongDescription" json:"longDescription"`
		ShortDescription   string         `db:"ShortDescription" json:"shortDescription"`
		HowToConsume       string         `db:"HowToConsume" json:"howToConsume"`
		ImageList          []ImageLists   `json:"imageList"`
		Price              float32        `db:"Price" json:"price"`
		FractionatedPrice  float32        `db:"fractionatedPrice" json:"fractionatedPrice"`
		Prescription       string         `db:"Prescription" json:"prescription"`
		Precautions        string         `db:"Precautions" json:"precautions"`
		SideEffects        string         `db:"SideEffects" json:"sideEffects"`
		Favorite           bool           `json:"favorite"`
		Presentation       string         `db:"Presentation" json:"presentation"`
		UnitQuantity       int16          `db:"UnitQuantity" json:"quantityUnits"`
		IsfractionalMode   bool           `db:"IsfractionalMode" json:"fractionalMode"`
		ActivePrinciples   []string       `json:"activePrinciples"`
		FractionatedForm   *string        `db:"fractionatedForm" json:"fractionatedForm"`
		FractionatedText   *string        `db:"fractionatedText" json:"fractionatedText"`
		UnFractionatedText *string        `db:"UnFractionatedText" json:"noFractionatedText"`
		StatusID           int            `db:"StatusID" json:"productStatusId"`
		ProductStatus      *string        `db:"productStatus" json:"productStatus"`
		MaxUnitSale        int32          `db:"maxUnitSale" json:"maxUnitSale"`
		CategoryList       []CategoryList `json:"categoryList"`
		Stock              int32          `db:"Stock" json:"stock"`
		FractionalStock    int32          `db:"FractionalStock" json:"fractionalStock"`
		ShowStockAlert     string         `db:"ShowStockAlert" json:"showStockAlert"`
	}

	type Level3 struct {
		Product ReturnQuery `json:"product"`
		// Product string `json:"product"`
	}

	type Level2 struct {
		Title    string   `json:"title"`
		Products []Level3 `json:"products"`
	}
	type Level1 struct {
		Type            string `json:"type"`
		BackgroundColor string `json:"backgroundColor"`
		TitleColor      string `json:"titleColor"`
		List            Level2 `json:"list"`
	}
	var level1s []Level1
	var level1 Level1
	for _, productRelatedRow := range productRelateds {
		log.Println(productRelatedRow)
		level1.Type = strconv.FormatInt(productRelatedRow.ProductRelatedTypeid, 10)

		type ProductRelatedType struct {
			ID              int64  `db:"id" json:"id"`
			Name            string `db:"name" json:"title"`
			BackgroundColor string `db:"backgroundColor" json:"backgroundColor"`
			TitleColor      string `db:"TitleColor" json:"titleColor"`
			IsVisible       string `db:"isVisible" json:"isVisible"`
		}
		query =
			`
			select prt.id,
				prt.name,
				prt.backgroundColor,
				prt.titleColor,
				prt.isVisible
		   from ProductRelatedType prt 
		 where prt.id = ?`
		var productRelatedType ProductRelatedType
		err = utils.DbmapMySQLInkafarma.SelectOne(&productRelatedType, query, productRelatedRow.ProductRelatedTypeid)
		if err != nil {
			log.Println(query)
			log.Println(err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if productRelatedType.IsVisible == "Y" {
			query =
				`  select p.id as id, 
			p.name as name, 
			p.longDescription as longDescription, 
			p.shortDescription as shortDescription, 
			p.howToConsume as howToConsume,
			p.price as price, 
			Case pp.unitQuantity
				when 0 then p.price
				when  null then p.price
				else p.price /pp.unitQuantity
			End as FractionatedPrice ,
			CASE p.prescriptionType
			WHEN 'RR' THEN "Retener Receta"
			WHEN 'PR' THEN "Presenta Receta"
			WHEN 'VL' THEN "Venta Libre"
			WHEN 'NA' THEN "No Aplica"
			ELSE "No Aplica."
			END as prescription,
			p.precautions as precautions, 
			p.sideEffects as sideEffects,
			p.statusId as statusId,
			p.stock as stock,
			( p.stock * pp.unitQuantity ) as fractionalStock,
			'N' as showStockAlert,
			pp.presentation,
			pp.unitQuantity , 
			pp.isfractionalMode,  
			unFractionatedText,
			pp.fractionatedForm,
			-- pp.pharmaceuticalForm,
			p.maxUnitSale,
			ps.name as productStatus,
			IFNULL(pp.fractionatedText," ") as fractionatedText
		from Product p 
			-- inner join productcategory pc on p.id=pc.productId 
			-- inner join Category c on pc.categoryId=c.id 
			inner join ProductPresentation pp on p.id = pp.productId
			inner join productstatus ps on ps.id = p.statusId
		where p.id in (` + productRelatedRow.ProductRelated + `)
		and (p.statusId in (1 , 3)) 
		order by p.ranking asc --  limit ?
		`

			var returnQuery []ReturnQuery
			_, err = utils.DbmapMySQLInkafarma.Select(&returnQuery, query)
			if err != nil {
				log.Println("Error ")
				log.Println(query)
				log.Println(err.Error())
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			var returnQueryWithImages []ReturnQuery
			for _, element := range returnQuery {

				/* ini Category List */
				var categoryList []CategoryList
				query = `
			select pc.categoryID as id,c.name
			from Product p
			 inner join productcategory pc on p.id=pc.productId
			 inner join Category c on pc.categoryId=c.id
			 where p.id = ?		
			`
				_, err := utils.DbmapMySQLInkafarma.Select(&categoryList, query, element.ID)
				if err != nil {
					log.Println(query)
					log.Println(err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
				}
				element.CategoryList = categoryList
				/* fin Category List */

				var productActiveComponent []ProductActiveComponent
				var activePrinciples []string
				query = `
			select name , concentration
			from ProductActiveComponent pac
				 join ActiveComponent ac on pac.activeComponentId = ac.id
			where productId = ?`
				_, err = utils.DbmapMySQLInkafarma.Select(&productActiveComponent, query, element.ID)
				if err != nil {
					log.Println(query)
					log.Println(err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
				}
				for _, ap := range productActiveComponent {
					activePrinciples = append(activePrinciples, ap.Name+" "+ap.Concentration)
				}
				element.ActivePrinciples = activePrinciples

				var imageLists []ImageLists
				query = `select imagePath from productimage pi where pi.productId = ?`
				_, err = utils.DbmapMySQLInkafarma.Select(&imageLists, query, element.ID)
				if err == nil {
					element.ImageList = imageLists
				} else {
					c.JSON(500, gin.H{"error": err.Error()})
				}
				// element.FractionatedPrice = element.Price / element.UnitQuantity
				returnQueryWithImages = append(returnQueryWithImages, element)

			}
			var level3s []Level3
			var level3 Level3
			for _, row := range returnQueryWithImages {
				level3.Product = row
				level3s = append(level3s, level3)
			}

			level1 := Level1{Type: strconv.FormatInt(productRelatedType.ID, 10),
				BackgroundColor: productRelatedType.BackgroundColor,
				TitleColor:      productRelatedType.TitleColor,
				List:            Level2{Title: productRelatedType.Name, Products: level3s}}
			level1s = append(level1s, level1)
		}

	}

	log.Println("************** STOP [ /getrelproducts (GetRelProducts) ] **************")
	if err == nil {
		c.JSON(200, gin.H{"itemList": level1s})
		// c.JSON(200, gin.H{"itemList": level1s})
	} else {
		log.Println("Error al final")
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}
}
