# Designer

A terminal design tool that allows for gradient art to be applied to text. It has full UTF-8 support, as well as a GIF to Text converter with different output options.

## 👨🏻‍🏫 Usage

To use Designer, all you need to do is compile the application and then execute it.

 1. `go mod tidy`
 2. `go build .`
 3. Execute the compiled binary.

## 🌈 Gradient

### 📄 Features

 - Multi-directional gradient support.
 - Full use of both RGB and HEX colours.
 - UTF-8 character support.
 - Unlimited colour input.

### 👨🏻‍🏫 Usage

Inside of your input file for gradients, you can use 'gradient tags.' These tags allow for you to choose the direction of a gradient, or wether you want the gradient to conclude at a certain point.

These tags are;
 - `diagonal` - Renders a diagonal gradient.
 - `radial` - Renders a radial gradient.
 - `vertical` - Renders a top to bottom gradient.
 - `horizontal` - Renders a left to right gradient.
 - `/` or `stop` - Stops the gradient from this point.

#### ⌨️ Example

```
<radial>All of this text from here to the next tag
will be rendered as a radial gradient.
<horizontal>This text will have a horizontal gradient applied to it.
So will this, </>but this text will revert back to having no gradient.
<vertical>This text
will have
a
vertical
gradient applied
to it.<stop>And this will have no colour.
```

## 📽️ GIF Converter

### 📄 Features

 - ANSI and ASCII GIF outputs.
 - Can be converted into Go files, Python files, and TermFX files.
 - Can fully resize GIF outputs to suit terminal size.
 - Can use both URLs and files when converting GIFs.
 - Frame-skipping (Makes playback less laggy on lower-powered machines.)
 - Ability to loop a GIF multiple times.
 - Supports `.gif`, `.jpeg`, `.png`, and `.webp` files.
    - Note: When rendering images, there will be no loop options as it is a static image.