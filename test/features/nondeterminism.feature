Feature: Nondeterminism
  Scenario: Apply rally function to a infinite list
    Given a file named "main.tisp" with:
    """
    (let (f) (prepend 42 (f)))

    (let many42 (f))
    (write (first many42))
    (let many42 (rest many42))
    (write (first many42))
    (let many42 (rest many42))
    (write (first many42))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    42
    42
    42
    """
