export function useImageHelper() {
  function uploadImage(model) {
    let form = new FormData();
    form.append("imageFile", model.imageFile[0]);
    form.append("scale", model.scale);
    form.append("noise", model.noise);
    form.append("uuid", model.uuid);
    return fetch("http://localhost:8080/api/v1/upload", {
      method: "POST",
      body: form,
    })
  }

  async function downloadImage(filename,status) {
    console.log(filename,status)
    if(status !== "done"){
      return
    }
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

  return {
    uploadImage,
    downloadImage,
  };
}
