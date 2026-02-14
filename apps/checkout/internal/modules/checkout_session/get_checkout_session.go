package checkout_session

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetCheckoutSession(c *routerx.Context, checkoutSession *models.CheckoutSession) {

	if err := c.DB().FindForID(checkoutSession.UserID, &checkoutSession.User); err != nil {
		c.ServerError("internal server error", err)
		return
	}
	if err := c.DB().FindForID(checkoutSession.PlanID, &checkoutSession.Plan); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	response := new(apix.CheckoutSession).Set(checkoutSession)
	response.ExternalSessionID = ""

	if checkoutSession.BillingAddressID != nil {
		var billingAddress models.Address
		if err := c.DB().FindForID(*checkoutSession.BillingAddressID, &billingAddress); err != nil {
			c.ServerError("internal server error", err)
			return
		}
		response.BillingAddress = &apix.Address{}
		response.BillingAddress.Set(&billingAddress)
	}

	// Company address is optional - only load if set
	if checkoutSession.CompanyAddressID != nil && *checkoutSession.CompanyAddressID > 0 {
		var companyAddress models.Address
		if err := c.DB().FindForID(*checkoutSession.CompanyAddressID, &companyAddress); err != nil {
			c.ServerError("internal server error", err)
			return
		}
		response.CompanyAddress = &apix.Address{}
		response.CompanyAddress.Set(&companyAddress)
	}

	// company - optional
	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", checkoutSession.AppID); err != nil {
		// Company is optional, don't fail
	} else {
		response.Company = &apix.Company{}
		response.Company.Set(&company)
	}

	plan := &checkoutSession.Plan
	response.Plan = &apix.Plan{}
	response.Plan.Set(plan)

	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	for _, item := range planItems {
		planItem := apix.PlanItem{}
		planItem.Set(&item)
		response.Plan.Items = append(response.Plan.Items, planItem)
	}

	if err := loadPlanFeatures(c.DB(), plan.ID, planItems, response.Plan); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	c.OK(response)
}
