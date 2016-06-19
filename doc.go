/*
Bayes is a simple package that compute bayesian calculation.

For example, it can classify emails in spam and non-spam categories using
some already classified messages.

Example:

	spam := NewCategory("spam")
	spam.Train("Content of spam message")
	spam.Train("Other content of spam message")

	nonspam := NewCategory("non-spam")
	nonspam.Train("Content without spam content")
	nonspam.Train("And another one")

	message := "Content of a spam message"
	res := Bayes(message, spam, []Category{spam, nonspam})

	// res gives a value between 0.0 and 1.0. 1.0 means that the
	// given content is a "spam"


You may use many Categories. The more you train them, the more you will
have better result.

Licence: MIT

*/
package bayes
