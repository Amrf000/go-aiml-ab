package ab

import (
	"fmt"
	"strings"
)

type TripleStore struct {
	IdCnt            int
	Name             string
	ChatSession      *Chat
	Bot              *Bot
	IdTriple         map[string]*Triple
	TripleStringId   map[string]string
	SubjectTriples   map[string]map[string]struct{}
	PredicateTriples map[string]map[string]struct{}
	ObjectTriples    map[string]map[string]struct{}
}

type Triple struct {
	Id        string
	Subject   string
	Predicate string
	Object    string
}

func NewTripleStore(name string, chatSession *Chat) *TripleStore {
	return &TripleStore{
		Name:             name,
		ChatSession:      chatSession,
		Bot:              chatSession.Bot,
		IdTriple:         make(map[string]*Triple),
		TripleStringId:   make(map[string]string),
		SubjectTriples:   make(map[string]map[string]struct{}),
		PredicateTriples: make(map[string]map[string]struct{}),
		ObjectTriples:    make(map[string]map[string]struct{}),
	}
}

func (ts *TripleStore) NewTriple(s, p, o string) *Triple {
	if ts.Bot != nil {
		s = ts.Bot.PreProcessor.Normalize(s)
		p = ts.Bot.PreProcessor.Normalize(p)
		o = ts.Bot.PreProcessor.Normalize(o)
	}

	if s != "" && p != "" && o != "" {
		triple := &Triple{
			Id:        fmt.Sprintf("%s%d", ts.Name, ts.IdCnt),
			Subject:   s,
			Predicate: p,
			Object:    o,
		}
		ts.IdCnt++
		return triple
	}
	return nil
}

func (ts *TripleStore) MapTriple(triple *Triple) string {
	id := triple.Id
	ts.IdTriple[id] = triple

	s := strings.ToUpper(triple.Subject)
	p := strings.ToUpper(triple.Predicate)
	o := strings.ToUpper(triple.Object)

	tripleString := fmt.Sprintf("%s:%s:%s", s, p, o)
	tripleString = strings.ToUpper(tripleString)

	if existingId, ok := ts.TripleStringId[tripleString]; ok {
		return existingId
	}

	ts.TripleStringId[tripleString] = id

	if _, ok := ts.SubjectTriples[s]; !ok {
		ts.SubjectTriples[s] = make(map[string]struct{})
	}
	ts.SubjectTriples[s][id] = struct{}{}

	if _, ok := ts.PredicateTriples[p]; !ok {
		ts.PredicateTriples[p] = make(map[string]struct{})
	}
	ts.PredicateTriples[p][id] = struct{}{}

	if _, ok := ts.ObjectTriples[o]; !ok {
		ts.ObjectTriples[o] = make(map[string]struct{})
	}
	ts.ObjectTriples[o][id] = struct{}{}

	return id
}

func (ts *TripleStore) UnMapTriple(triple *Triple) string {
	id := UndefinedTriple

	s := strings.ToUpper(triple.Subject)
	p := strings.ToUpper(triple.Predicate)
	o := strings.ToUpper(triple.Object)

	tripleString := fmt.Sprintf("%s:%s:%s", s, p, o)
	fmt.Println("unMapTriple", tripleString)
	tripleString = strings.ToUpper(tripleString)

	if existingId, ok := ts.TripleStringId[tripleString]; ok {
		triple = ts.IdTriple[existingId]
		fmt.Println("unMapTriple", triple)
		if triple != nil {
			id = triple.Id
			delete(ts.IdTriple, id)
			delete(ts.TripleStringId, tripleString)

			delete(ts.SubjectTriples[s], id)
			delete(ts.PredicateTriples[p], id)
			delete(ts.ObjectTriples[o], id)
		}
	}

	return id
}

func (ts *TripleStore) AllTriples() map[string]*Triple {
	result := make(map[string]*Triple, len(ts.IdTriple))
	for id, triple := range ts.IdTriple {
		result[id] = triple
	}
	return result
}

func (ts *TripleStore) AddTriple(subject, predicate, object string) string {
	if subject == "" || predicate == "" || object == "" {
		return UndefinedTriple
	}

	triple := ts.NewTriple(subject, predicate, object)
	if triple == nil {
		return UndefinedTriple
	}

	id := ts.MapTriple(triple)
	return id
}

func (ts *TripleStore) DeleteTriple(subject, predicate, object string) string {
	if subject == "" || predicate == "" || object == "" {
		return UndefinedTriple
	}

	triple := ts.NewTriple(subject, predicate, object)
	if triple == nil {
		return UndefinedTriple
	}

	id := ts.UnMapTriple(triple)
	return id
}

func (ts *TripleStore) PrintTriples() {
	for id, triple := range ts.IdTriple {
		fmt.Printf("%s:%s:%s:%s\n", id, triple.Subject, triple.Predicate, triple.Object)
	}
}

func (ts *TripleStore) GetTriples(s, p, o string) map[string]*Triple {
	subjectSet := ts.AllTriples()
	predicateSet := ts.AllTriples()
	objectSet := ts.AllTriples()
	resultSet := make(map[string]*Triple)

	fmt.Println("TripleStore: getTriples [", len(ts.IdTriple), "] ", s, ":", p, ":", o)

	if s == "" || strings.HasPrefix(s, "?") {
		for id := range subjectSet {
			resultSet[id] = ts.IdTriple[id]
		}
	} else {
		s = strings.ToUpper(s)
		if subjTriples, ok := ts.SubjectTriples[s]; ok {
			for id := range subjTriples {
				resultSet[id] = ts.IdTriple[id]
			}
		}
	}

	if p == "" || strings.HasPrefix(p, "?") {
		for id := range predicateSet {
			resultSet[id] = ts.IdTriple[id]
		}
	} else {
		p = strings.ToUpper(p)
		if predTriples, ok := ts.PredicateTriples[p]; ok {
			for id := range predTriples {
				resultSet[id] = ts.IdTriple[id]
			}
		}
	}

	if o == "" || strings.HasPrefix(o, "?") {
		for id := range objectSet {
			resultSet[id] = ts.IdTriple[id]
		}
	} else {
		o = strings.ToUpper(o)
		if objTriples, ok := ts.ObjectTriples[o]; ok {
			for id := range objTriples {
				resultSet[id] = ts.IdTriple[id]
			}
		}
	}

	return resultSet
}

func (ts *TripleStore) GetSubjects(triples map[string]*Triple) map[string]struct{} {
	resultSet := make(map[string]struct{})
	for _, triple := range triples {
		resultSet[triple.Subject] = struct{}{}
	}
	return resultSet
}

func (ts *TripleStore) GetPredicates(triples map[string]*Triple) map[string]struct{} {
	resultSet := make(map[string]struct{})
	for _, triple := range triples {
		resultSet[triple.Predicate] = struct{}{}
	}
	return resultSet
}

func (ts *TripleStore) GetObjects(triples map[string]*Triple) map[string]struct{} {
	resultSet := make(map[string]struct{})
	for _, triple := range triples {
		resultSet[triple.Object] = struct{}{}
	}
	return resultSet
}

func (ts *TripleStore) FormatAIMLTripleList(triples map[string]*Triple) string {
	var result strings.Builder
	result.WriteString(DefaultListItem)
	for id := range triples {
		result.WriteString(fmt.Sprintf("%s ", id))
	}
	return strings.TrimSpace(result.String())
}

func (ts *TripleStore) GetSubject(id string) string {
	if triple, ok := ts.IdTriple[id]; ok {
		return triple.Subject
	}
	return "Unknown subject"
}

func (ts *TripleStore) GetPredicate(id string) string {
	if triple, ok := ts.IdTriple[id]; ok {
		return triple.Predicate
	}
	return "Unknown predicate"
}

func (ts *TripleStore) GetObject(id string) string {
	if triple, ok := ts.IdTriple[id]; ok {
		return triple.Object
	}
	return "Unknown object"
}

func (ts *TripleStore) StringTriple(id string) string {
	if triple, ok := ts.IdTriple[id]; ok {
		return fmt.Sprintf("%s %s %s %s", id, triple.Subject, triple.Predicate, triple.Object)
	}
	return ""
}

func (ts *TripleStore) PrintAllTriples() {
	for id := range ts.IdTriple {
		fmt.Println(ts.StringTriple(id))
	}
}

func (ts *TripleStore) Select(vars, visibleVars map[string]bool, clauses []*Clause) []*Tuple {
	result := []*Tuple{}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Something went wrong with select ", visibleVars)
			fmt.Println(r)
		}
	}()

	tuple := NewTuple(vars, visibleVars, nil)
	result = ts.SelectFromRemainingClauses(tuple, clauses)

	if TraceMode {
		for _, t := range result {
			fmt.Println(t.PrintTuple())
		}
	}

	return result
}

func (ts *TripleStore) AdjustClause(tuple *Tuple, clause *Clause) *Clause {
	vars := tuple.GetVars()
	subj := clause.Subj
	pred := clause.Pred
	obj := clause.Obj

	newClause := &Clause{Subj: subj, Pred: pred, Obj: obj}

	if _, ok := vars[subj]; ok {
		value := tuple.GetValue(subj)
		if value != UnboundVariable {
			newClause.Subj = value
		}
	}

	if _, ok := vars[pred]; ok {
		value := tuple.GetValue(pred)
		if value != UnboundVariable {
			newClause.Pred = value
		}
	}

	if _, ok := vars[obj]; ok {
		value := tuple.GetValue(obj)
		if value != UnboundVariable {
			newClause.Obj = value
		}
	}

	return newClause
}

func (ts *TripleStore) BindTuple(partial *Tuple, triple string, clause *Clause) *Tuple {
	tuple := NewTupleClone(partial)
	if strings.HasPrefix(clause.Subj, "?") {
		tuple.Bind(clause.Subj, ts.GetSubject(triple))
	}
	if strings.HasPrefix(clause.Pred, "?") {
		tuple.Bind(clause.Pred, ts.GetPredicate(triple))
	}
	if strings.HasPrefix(clause.Obj, "?") {
		tuple.Bind(clause.Obj, ts.GetObject(triple))
	}
	return tuple
}

func removeDuplicate(sliceList []*Tuple) []*Tuple {
	allKeys := make(map[int]*Tuple)
	list := []*Tuple{}
	for _, item := range sliceList {
		if _, value := allKeys[item.HashCode()]; !value {
			allKeys[item.HashCode()] = item
			list = append(list, item)
		}
	}
	return list
}

func (ts *TripleStore) SelectFromSingleClause(partial *Tuple, clause *Clause, affirm bool) []*Tuple {
	result := []*Tuple{}
	triples := ts.GetTriples(clause.Subj, clause.Pred, clause.Obj)

	if affirm {
		for id := range triples {
			tuple := ts.BindTuple(partial, id, clause)
			result = append(result, tuple)
		}
	} else {
		if len(triples) == 0 {
			result = append(result, partial)
		}
	}

	result = removeDuplicate(result)

	return result
}

func (ts *TripleStore) SelectFromRemainingClauses(partial *Tuple, clauses []*Clause) []*Tuple {
	result := []*Tuple{}

	clause := clauses[0]
	clause = ts.AdjustClause(partial, clause)
	tuples := ts.SelectFromSingleClause(partial, clause, clause.Affirm)

	if len(clauses) > 1 {
		remainingClauses := make([]*Clause, len(clauses)-1)
		copy(remainingClauses, clauses[1:])
		for _, tuple := range tuples {
			result = append(result, ts.SelectFromRemainingClauses(tuple, remainingClauses)...)
		}
	} else {
		result = tuples
	}

	return result
}
