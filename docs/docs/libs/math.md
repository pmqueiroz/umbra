# Math

Perform mathematical operations such as absolute value, square root, trigonometric calculations, and power functions.

## Functions

### `abs(x num) num`
Returns the absolute value of a number.

- **Parameters:**
  - `x` (`num`): The number to get the absolute value of.
- **Returns:**
  - (`num`): The absolute value of the number.

### `sqrt(x num) num`
Calculates the square root of a number using the Newton-Raphson method.

- **Parameters:**
  - `x` (`num`): The number to calculate the square root of.
- **Returns:**
  - (`num`): The square root of the number, or `NaN` if the number is negative.

### `floor(x num) num`
Rounds a number down to the nearest integer, handling negative numbers correctly.

- **Parameters:**
  - `x` (`num`): The number to be rounded down.
- **Returns:**
  - (`num`): The largest integer less than or equal to the number.

### `ceil(x num) num`
Rounds a number up to the nearest integer, handling positive numbers correctly.

- **Parameters:**
  - `x` (`num`): The number to be rounded up.
- **Returns:**
  - (`num`): The smallest integer greater than or equal to the number.

### `pow(base num, exponent num) num`
Calculates the power of a number.

- **Parameters:**
  - `base` (`num`): The base number.
  - `exponent` (`num`): The exponent to raise the base to.
- **Returns:**
  - (`num`): The result of raising the base to the given exponent.

### `sin(x num) num`
Calculates the sine of a number using the Taylor series expansion.

- **Parameters:**
  - `x` (`num`): The number (in radians) to calculate the sine of.
- **Returns:**
  - (`num`): The sine of the number.

### `cos(x num) num`
Calculates the cosine of a number using the Taylor series expansion.

- **Parameters:**
  - `x` (`num`): The number (in radians) to calculate the cosine of.
- **Returns:**
  - (`num`): The cosine of the number.

### `tan(x num) num`
Calculates the tangent of a number.

- **Parameters:**
  - `x` (`num`): The number (in radians) to calculate the tangent of.
- **Returns:**
  - (`num`): The tangent of the number.
