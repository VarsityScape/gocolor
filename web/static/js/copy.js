async function copy(x) {
    await navigator.clipboard.writeText(x);
    alert("Copied the text: " + x);
  }