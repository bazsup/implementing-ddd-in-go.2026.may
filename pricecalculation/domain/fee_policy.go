package domain

const numberOfVisitsForFeeForPrivateVisitors = uint(3)
const additionalFeeOn3rdVisitPerMonthForPrivateVisitorsPercentage = uint(5)

type Visits interface {
	NumberOfVisitsInMonthOfLastVisit() int
}

type FeePolicy struct {
	visitor ExternalVisitor
}

func NewFeePolicy(visitor ExternalVisitor) FeePolicy {
	return FeePolicy{visitor: visitor}
}

func (p FeePolicy) AddFee(visits Visits, totalPrice Price) Price {
	switch p.visitor.Type() {
	case ExternalVisitorTypePrivate:
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
