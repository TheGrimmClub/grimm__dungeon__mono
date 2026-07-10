"""grimm — the teaching package, seeded into the dungeon work dir.

This is the Actor-focused subset of TheGrimmClub's `grimm__python__zero` package
(the source of truth: https://github.com/TheGrimmClub/grimm__python__zero). It is
written into ~/.grimm/work so behavioral puzzle solutions can:

    from grimm import Actor
    print(Actor(name="dein-name"))
"""

from .actor import Actor

__all__ = ["Actor"]
__version__ = "0.1.0"
