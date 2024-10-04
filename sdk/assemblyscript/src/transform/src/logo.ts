import { Colors } from "assemblyscript/util/terminal.js";
import { WriteStream as FSWriteStream } from "fs";
import { WriteStream as TTYWriteStream } from "tty";

export default (stream: FSWriteStream | TTYWriteStream, svg = false) => {
  if (svg) {
    writeMarkdownLogo(stream);
  } else {
    writeAsciiLogo(stream);
  }
};

function writeAsciiLogo(stream: FSWriteStream | TTYWriteStream) {
  const logo = String.raw`
     __ __                                __   
    / // /_ _____  ___ ______ _  ___  ___/ /__ 
   / _  / // / _ \/ -_) __/  ' \/ _ \/ _  / -_)
  /_//_/\_, / .__/\__/_/ /_/_/_/\___/\_,_/\__/ 
       /___/_/                                 
`;
  const colors = new Colors(stream as { isTTY: boolean });
  stream.write(colors.blue(logo) + "\n");
}

function writeMarkdownLogo(stream: FSWriteStream | TTYWriteStream) {
  const logo = String.raw`
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/hypermodeinc/.github/main/images/hypermode-white.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/hypermodeinc/.github/main/images/hypermode-black.svg">
  <img alt="Hypermode" src="https://raw.githubusercontent.com/hypermodeinc/.github/main/images/hypermode-black.svg">
</picture>
`;
  stream.write(logo + "\n");
}
