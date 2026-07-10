"""The Actor — your very first Python class.

A *class* is a blueprint. `Actor` is the blueprint; each `Actor()` you make is
one character (an "instance") in the Grimm world. Read this file top to bottom —
it introduces the day-two ideas one at a time.
"""

# A default so that `Actor()` (no name) still works.
DEFAULT_NAME = "Namenloser"  # German for "the nameless one"


class Actor:
    """A character in the Grimm world — a player's avatar.

    Try it:

        from grimm import Actor

        me = Actor()
        print(me)
    """

    def __init__(self, name=DEFAULT_NAME):
        # __init__ runs automatically when you write `Actor(...)`.
        # `self` is *this* actor — the one being created right now.
        # We store the name on the actor so it remembers it later.
        self._name = name
        self._in_dungeon = False

    def name(self):
        """Give back this actor's name.

        It's a *method* (a function that belongs to the actor), so you call it
        with parentheses: `me.name()`.
        """
        return self._name

    def enter_dungeon(self, launch=False):
        """Step through the gate into the Grimm dungeon.

        By default this just tells the story. Pass ``launch=True`` to actually
        start the real adventure (the `grimm` console program):

            me.enter_dungeon(launch=True)
        """
        self._in_dungeon = True
        # An f-string ("formatted string") drops the name straight into the text.
        print(f"{self._name} betritt das Verlies. Die Tür fällt ins Schloss. \N{CANDLE}")
        if launch:
            # Import here so the toy Actor doesn't drag in the launcher unless asked.
            from .dungeon import Dungeon

            return Dungeon().enter()
        return None

    def __str__(self):
        # __str__ decides what `print(actor)` shows. Without it, printing an
        # actor would show something ugly like <grimm.actor.Actor object at 0x...>.
        ort = "im Verlies" if self._in_dungeon else "vor dem Verlies"
        return f"\N{STANDING PERSON} {self._name} wartet {ort}."

    def __repr__(self):
        # __repr__ is the "developer" view, handy in the REPL and in tests.
        return f"Actor(name={self._name!r})"
