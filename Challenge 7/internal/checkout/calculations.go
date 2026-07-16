package checkout

import "time"

/**
 * Returns check out date + 14 days
 */
func WhenDue(checkoutDate time.Time) time.Time {
	return checkoutDate.Add(time.Duration(14 * 24 * time.Hour))
}

/**
 * Check if the book is late
 */
func IsLate(dueDate time.Time) bool {
	if time.Now().After(dueDate) {
		return true
	}
	return false
}

/**
 * If late, how much is owed
 * Simple calc 25 cents per day late
 * If days difference is <= 1, 25 cents
 */
func AmountDue(dueDate time.Time) float64 {
	if IsLate(dueDate) {
		diff := time.Now().Sub(dueDate)
		days := diff.Hours() / 24
		if days < 1 {
			return 0.25
		} else {
			return days * 0.25
		}
	}
	return 0.00
}
