export default async function* lineReader(reader: ReadableStreamDefaultReader<Uint8Array>) {
    const utf8Decoder = new TextDecoder("utf-8");
   
    let {value: chunkVal, done: readerDone} = await reader.read();
    let chunk = chunkVal ? utf8Decoder.decode(chunkVal, {stream: true}) : "";
  
    let re = /\r\n|\n|\r/gm;
    let startIndex = 0;
  
    for (;;) {
      let result = re.exec(chunk);
      if (!result) {
        if (readerDone) {
          break;
        }
        let remainder = chunk.substr(startIndex);
        ({value: chunkVal, done: readerDone} = await reader.read());
        chunk = remainder + (chunkVal ? utf8Decoder.decode(chunkVal, {stream: true}) : "");
        startIndex = re.lastIndex = 0;
        continue;
      }
      yield chunk.substring(startIndex, result.index);
      startIndex = re.lastIndex;
    }
    if (startIndex < chunk.length) {
      // last line didn't end in a newline char
      yield chunk.substr(startIndex);
    }
  }
