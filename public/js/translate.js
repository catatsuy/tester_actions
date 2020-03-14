window.addEventListener('DOMContentLoaded', (event) => {
  const input = document.querySelector('#input');
  const target = document.querySelector('#target');
  document.querySelector('#translate').addEventListener('click', (event) => {
    fetch('/api/translate', {
      method: "POST",
      headers: {
        "Content-Type": "application/json; charset=utf-8",
      },
      body: JSON.stringify({input: input.value}),
    })
      .then((response) => {
        return response.json();
      })
      .then((json) => {
        target.value = json.output;
        navigator.clipboard.writeText(json.output);
      })
  });
});
