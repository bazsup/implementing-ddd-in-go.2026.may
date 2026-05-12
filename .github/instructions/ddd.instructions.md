---
applyTo: "pricecalculation/**"
---

# Tactical DDD

Design, refactor, analyze, and review code by applying the principles and patterns of tactical domain-driven design.

## Principles

1. **Isolate domain logic**
2. **Use rich domain language**
3. **Orchestrate with applications**
4. **Avoid anemic domain model**
5. **Separate generic concepts**
6. **Make the implicit explicit... like your life depends on it**
7. **Design aggregates around invariants**
8. **Extract immutable value objects liberally**
9. **Repositories are for loading and saving full aggregates**

---

## 1. Isolate domain logic

**What:** Domain logic is not mixed with technical code like HTTP and database transactions.

**Why:** Easier to understand the most important part of the code, easier to validate with domain experts, easier to test and evolve, easier to plan and implement new features.

**Test:** Could a domain expert read the code? Can the code be unit tested without mocks or spinning up databases?

---

## 2. Use rich domain language

**What:** Names in code match exactly what domain experts say. No programmer jargon. No generic names.

**Why:** Translation between code-speak and business-speak causes bugs.

**Test:** Would a domain expert recognize this name? If you'd need to translate it for them, it's wrong.

**Common generic terms to watch for:**
- `Manager`, `Handler`, `Processor`, `Helper`, `Util`
- `Data`, `Info`, `Item` (when domain terms exist)
- `process`, `handle`, `execute` (what does it actually DO?)

---

## 3. Orchestrate with applications

**What:** An application is a user goal—something a user would recognize as an action they can perform.

**Why:** Applications define the entry points to your domain.

**Test (the menu test):** If you described your application's features to a user like a menu, would this be on it?

---

## 4. Avoid anemic domain model

**What:** Domain logic lives in domain objects, not in applications. Applications orchestrate; domain objects decide.

**Why:** When business rules leak into applications, they scatter, duplicate, and diverge.

**Test:** Is your application making business decisions, or just coordinating? If the application contains if/else business logic, you likely have an anemic model.

---

## 5. Separate generic concepts

**What:** Generic capabilities that aren't specific to your domain live separately from domain-specific logic.

**Why:** A retry mechanism, a caching layer — these aren't YOUR domain. Mixing them obscures what's specific to your business.

**Test:** Would this code exist in a completely different business domain? If yes, it's generic.

---

## 6. Make the implicit explicit... like your life depends on it

**What:** Strive for maximum expressiveness. Go as far as possible to identify and name domain concepts in code.

**Why:** Maximum alignment optimizes communication between engineers and domain experts.

**Test:** Could you discuss this code with a domain expert without translation?

---

## 7. Design aggregates around invariants

**What:** An aggregate is a cluster of objects that must be consistent together. The aggregate root enforces the rules.

**Why:** Without clear boundaries, inconsistent states creep in.

**Test:** What must be true at all times? The objects involved in those rules form an aggregate.

**Aggregate rules:**
- One root entity per aggregate
- External code accesses only through the root
- The root enforces all invariants
- Reference other aggregates by ID, not object

---

## 8. Extract immutable value objects liberally

**What:** When something is defined by its attributes (not identity), make it an immutable value object.

**Why:** Value objects are simple, can't change unexpectedly, easy to test, make domain concepts explicit.

**Test:** Does this need a unique ID to track it over time? No? It's probably a value object.

**Good candidates:**
- Money, Currency, Percentage
- DateRange, TimeSlot, Duration
- Address, Coordinates, Distance
- EmailAddress, PhoneNumber, URL
- Quantity, Weight, Temperature

---

## 9. Repositories are for loading and saving full aggregates

The job of a repository is to load and save entire aggregates — not partial aggregates or nested entities.

A repository should not exist for a domain object that is not an aggregate root.

The `save` method takes the full aggregate. If you just want to query information to display without modifying state, create a separate read model object.

---

## 10. Naming conventions

**Value objects, entities, aggregates** — the domain noun. No pattern name appended.
```
// ❌ WRONG:  PriceValueObject   OrderAggregate
// ✅ RIGHT:  Price              Order
```

**Repositories** — plural noun for the interface; storage mechanism as prefix for implementations.
```
// ❌ WRONG:  VisitRepository
// ✅ RIGHT:  Visits (interface)   InMemoryVisits (implementation)
```

**Domain events** — past tense.
```
// ❌ WRONG:  PriceCalculatedEvent
// ✅ RIGHT:  PriceCalculated
```

**Bounded contexts** — gerund form.
```
// ❌ WRONG:  account/   invoice/
// ✅ RIGHT:  accounting/   invoicing/
```

**Domain services** — the domain concept performing the action. No "Service" suffix.
```
// ❌ WRONG:  PricingService
// ✅ RIGHT:  Pricer   ShippingCostCalculator
```

**Applications** — verb phrase describing the user action. No "Application" or "Service" suffix.

**Errors** — domain language describing what went wrong. No generic `Exception` or `NotFoundException`.
```
// ❌ WRONG:  NotFoundException
// ✅ RIGHT:  OrderNotFound   InsufficientFunds
```

---

## Mandatory Checklist

When designing, refactoring, analyzing, or reviewing code:

1. [ ] Verify domain is isolated from infrastructure (no DB/HTTP/logging in domain)
2. [ ] Verify names are from YOUR domain, not generic developer jargon
3. [ ] Verify applications are intentions of users (apply the menu test)
4. [ ] Verify business logic lives in domain objects, applications only orchestrate
5. [ ] Verify states are modeled as distinct types where appropriate
6. [ ] Verify hidden domain concepts are extracted and named explicitly
7. [ ] Verify aggregates are designed around invariants
8. [ ] Verify values are extracted into value objects expressing a domain concept
9. [ ] Verify no abuse of hydrate methods; each creation scenario has a dedicated factory method
10. [ ] Verify file and class names follow naming conventions: no pattern suffixes, repositories are plural, domain events are past tense, bounded contexts are gerunds
