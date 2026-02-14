package handlers

import (
	"net/http"

	"github.com/useportcall/portcall/apps/quote/internal/i18n"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetQuote(c *routerx.Context) {
	var quote models.Quote
	if err := c.DB().FindFirst(&quote, "public_id = ?", c.Param("id")); err != nil {
		c.NotFound("quote not found")
		return
	}
	if !verifyQuoteAccess(c, &quote) {
		return
	}

	var plan models.Plan
	if err := c.DB().FindForID(quote.PlanID, &plan); err != nil {
		c.NotFound("quote not found")
		return
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", quote.AppID); err != nil {
		c.NotFound("company not found")
		return
	}

	var companyAddress models.Address
	if err := c.DB().FindForID(company.BillingAddressID, &companyAddress); err != nil {
		c.NotFound("company address not found")
		return
	}

	var recipientAddress models.Address
	if err := c.DB().FindForID(quote.RecipientAddressID, &recipientAddress); err != nil {
		c.NotFound("recipient address not found")
		return
	}

	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.NotFound("quote not found")
		return
	}

	inst := i18n.GetInstance()
	lang := inst.GetLanguage(c.Request)

	itemDatas, total, basePrice := buildQuoteItems(planItems, lang, inst, plan.Interval)

	features, err := loadFeatureNames(c, plan.ID)
	if err != nil {
		c.ServerError("failed to load features", err)
		return
	}

	service := c.DefaultQuery("service", plan.Name)
	accessToken := c.Query("qt")

	data := assembleQuoteData(
		&quote, &plan, &company,
		&companyAddress, &recipientAddress,
		itemDatas, features, total, basePrice, accessToken,
		inst, lang,
	)
	data.Service = service

	c.HTML(http.StatusOK, "quote.html", data)
}
