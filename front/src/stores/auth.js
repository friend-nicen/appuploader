import {defineStore} from "pinia";

export default defineStore("auth",{
  state(){
      return {
          token:""
      }
  }
});