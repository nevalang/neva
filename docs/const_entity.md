# Constant Entity

Constant is an entity that consist of either _message_ or _entity reference_ to other constant. Message can include references to other constants. Constant messages can be infinitely nested. Constants may refer imported constants from other packages. _Components_ are only entities that can refer constants, that are not constants themselves - they can refer to constants via _compiler directives_ and from their _networks_.
