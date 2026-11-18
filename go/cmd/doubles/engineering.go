package main

type Engineering interface {
	DoEngineering(*Ship)
}

type NoEngineering struct{}

func (_ NoEngineering) DoEngineering(_ *Ship) {}

type MaxEngineering struct {
	BarrierTransferCost int
}

func (me MaxEngineering) DoEngineering(ship *Ship) {
	// Try shifting another barrier to cover the damaged barrier.
	const otherBarrierToMove = 5
	if ship.BarrierDamage >= otherBarrierToMove && ship.RemainingStress() >= me.BarrierTransferCost {
		for idx := range ship.OtherBarriersMax {
			if ship.RemainingOtherBarrier(idx) >= otherBarrierToMove {
				ship.BarrierDamage -= otherBarrierToMove
				ship.DamageOtherBarrier(idx, otherBarrierToMove)
				ship.Stress += me.BarrierTransferCost
				break
			}
		}
	}

	// Use barrier generator to refill barrier if able to.
	barrierBoost := func() {
		const boostBarrierBy = 2
		const barrierBoostStress = 1

		if ship.RemainingStress() < barrierBoostStress {
			// Ship can't afford the additional stress.
			return
		}

		if ship.BarrierDamage >= boostBarrierBy {
			// Boost the forward barrier if we can.
			ship.BarrierDamage -= boostBarrierBy
			ship.Stress += barrierBoostStress
			return
		}

		// Otherwise try to find another "other" barrier to boost.
		for idx := range ship.OtherBarriersMax {
			if ship.GetOtherBarrierDamage(idx) >= boostBarrierBy {
				ship.OtherBarriersDamage[idx] -= boostBarrierBy
				ship.Stress += barrierBoostStress
				return
			}
		}
	}
	if ship.BarrierBoost {
		barrierBoost()
		if ship.ImprovedBarrierGenerator {
			barrierBoost()
		}
	}
}
