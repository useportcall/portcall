package checkout_session

import (
	"fmt"
	"os"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetCheckoutSession(c *routerx.Context) {
	sessionID := c.Param("id")

	var checkoutSession models.CheckoutSession
	if err := c.DB().FindFirst(&checkoutSession, "public_id = ?", sessionID); err != nil {
		c.NotFound("checkout session not found")
		return
	}

	response := new(apix.CheckoutSession).Set(&checkoutSession)
	response.URL = fmt.Sprintf("%s/%s", os.Getenv("CHECKOUT_URL"), checkoutSession.PublicID)

	if checkoutSession.BillingAddressID != nil {
		var billingAddress models.Address
		if err := c.DB().FindForID(*checkoutSession.BillingAddressID, &billingAddress); err != nil {
			c.ServerError("internal server error", err)
			return
		}
		response.BillingAddress = &apix.Address{}
		response.BillingAddress.Set(&billingAddress)
	}

	var companyAddress models.Address
	if err := c.DB().FindForID(checkoutSession.CompanyAddressID, &companyAddress); err != nil {
		c.ServerError("internal server error", err)
		return
	}
	response.CompanyAddress = &apix.Address{}
	response.CompanyAddress.Set(&companyAddress)

	// company
	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", checkoutSession.AppID); err != nil {
		c.ServerError("internal server error", err)
		return
	}
	response.Company = &apix.Company{}
	response.Company.Set(&company)

	var dbPlan *models.Plan
	if err := c.DB().FindForID(checkoutSession.PlanID, &dbPlan); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	response.Plan = &apix.Plan{}
	response.Plan.Set(dbPlan)

	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", dbPlan.ID); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	for _, item := range planItems {
		planItem := apix.PlanItem{}
		planItem.Set(&item)
		response.Plan.Items = append(response.Plan.Items, planItem)
	}

	var planFeatures []models.PlanFeature
	if err := c.DB().List(&planFeatures, "plan_id = ?", dbPlan.ID); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	for _, pf := range planFeatures {
		var feature models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &feature); err != nil {
			c.ServerError("internal server error", err)
			return
		}

		if feature.IsMetered {
			res := apix.PlanFeature{}
			res.Set(&pf)
			res.Feature = apix.Feature{ID: feature.PublicID, IsMetered: feature.IsMetered}

			// find plan item
			planItem := apix.PlanItem{}
			for _, item := range planItems {
				if item.ID == pf.PlanItemID {
					planItem.Set(&item)
					res.PlanItem = &planItem
					break
				}
			}

			response.Plan.MeteredFeatures = append(response.Plan.MeteredFeatures, res)
		} else {
			res := apix.Feature{ID: feature.PublicID, IsMetered: feature.IsMetered}
			response.Plan.Features = append(response.Plan.Features, res)
		}
	}

	c.OK(response)
}
