package ab

type Clause struct {
	Subj   string
	Pred   string
	Obj    string
	Affirm bool
}

func NewClause(s, p, o string) *Clause {
	return &Clause{
		Subj:   s,
		Pred:   p,
		Obj:    o,
		Affirm: true,
	}
}

func NewClauseWithAffirm(s, p, o string, affirm bool) *Clause {
	return &Clause{
		Subj:   s,
		Pred:   p,
		Obj:    o,
		Affirm: affirm,
	}
}

func NewClauseFromClause(clause *Clause) *Clause {
	return &Clause{
		Subj:   clause.Subj,
		Pred:   clause.Pred,
		Obj:    clause.Obj,
		Affirm: clause.Affirm,
	}
}
