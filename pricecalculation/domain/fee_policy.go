package domain

const numberOfVisitsForFeeForPrivateVisitors = uint(3)
const additionalFeeOn3rdVisitPerMonthForPrivateVisitorsPercentage = uint(5)

type Visits interface {
	NumberOfVisitsInMonthOfLastVisit() int
}

type FeePolicy struct {
	customer Customer
}

func NewFeePolicy(customer Customer) FeePolicy {
	return FeePolicy{customer: customer}
}

func (p FeePolicy) AddFee(visits Visits, totalPrice Price) Price {
	switch p.customer.Type() {
	case customerTypePrivate:
		return p.addFeeForPrivateVisitor(visits, totalPrice)
	default:
		return totalPrice
	}
}

func (p FeePolicy) addFeeForPrivateVisitor(visits Visits, totalPrice Price) Price {
	if visits.NumberOfVisitsInMonthOfLastVisit() >= int(numberOfVisitsForFeeForPrivateVisitors) {
		totalPrice = totalPrice.AddFee(additionalFeeOn3rdVisitPerMonthForPrivateVisitorsPercentage)
	}
	return totalPrice
}
