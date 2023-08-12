export function useImageHelper() {
  function uploadImage(model) {
    let form = new FormData();
    form.append("imageFile", model.imageFile[0]);
    form.append("scale", model.scale);
    form.append("noise", model.noise);
    form.append("uuid", session.uuid);
    return fetch("http://localhost:8080/api/v1/upload", {
      method: "POST",
      body: form,
    })
  }

  async function downloadImage(filename) {
    console.log("Downloading", filename);
    const image = await fetch(
      "http://localhost:8080/api/v1/download-image?" +
        new URLSearchParams({
          filename: filename,
        })
    );
    const imageBlob = await image.blob();
    const imageURL = URL.createObjectURL(imageBlob);

    const link = document.createElement("a");
    link.href = imageURL;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  async function getImages() {
    const response = await fetch(
      "http://localhost:8080/api/v1/get-images?" +
        new URLSearchParams({
          uuid: session.uuid,
        })
    );
    const jsonData = await response.json();
    return jsonData;
  }
  return {
    uploadImage,
    downloadImage,
    getImages,
  };
}
