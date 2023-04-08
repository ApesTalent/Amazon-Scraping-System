import { toast } from "react-toastify";
import axios from "axios";
import cookie from "react-cookies";

export const API_URL = "http://localhost:8000/api";

export const Client = (auth = false) => {
  const defaultOptions = {
    headers: auth
      ? {
          Authorization: "Bearer " + cookie.load("token"),
        }
      : {},
  };

  return {
    get: (url, options = {}) =>
      axios.get(url, { ...defaultOptions, ...options }),
    post: (url, data, options = {}) =>
      axios.post(url, data, { ...defaultOptions, ...options }),
    put: (url, data, options = {}) =>
      axios.put(url, data, { ...defaultOptions, ...options }),
    delete: (url, options = {}) =>
      axios.delete(url, { ...defaultOptions, ...options }),
  };
};

export const createNotification = (type, message) => {
  toast.configure({
    position: "top-right",
    autoClose: 5000,
    hideProgressBar: false,
    closeOnClick: true,
    pauseOnHover: true,
    draggable: true,
  });
  switch (type) {
    case "success":
      toast.success(message);
      return;
    case "warn":
      toast.warn(message);
      return;
    case "error":
      toast.error(message);
      return;
    case "info":
      toast.info(message);
      return;
    default:
      toast(message);
  }
};

export function errorMessage(err) {
  if (!err.response) return err.message;
  if (err.response.data) {
    return err.response.data;
  }
  return err.message;
}
