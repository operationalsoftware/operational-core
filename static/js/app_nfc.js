/**
 * @namespace AppNFC
 */
const AppNFC = (() => {
  /**
   * @typedef {Object} NFCTagData
   * @property {string} tagContent - The content of the NFC tag.
   * @property {"text"|"url"} recordType - The type of record.
   */

  /**
   * Check if the device has NFC support.
   * @returns {Promise<boolean>}
   */
  async function checkDeviceHasNFC() {
    return "NDEFReader" in window;
  }

  /**
   * Read NFC tag data.
   * @returns {Promise<NFCTagData>}
   */
  async function readNFCTag() {
    try {
      return await Promise.race([
        new Promise((_, reject) => {
          setTimeout(() => {
            reject(new Error("operation timed out"));
          }, 3000);
        }),
        new Promise(async (resolve, reject) => {
          const ndef = new NDEFReader();
          await ndef.scan();

          ndef.onreading = (event) => {
            const record = event.message.records[0];
            const decoder = new TextDecoder(record.encoding);
            const tagContent = decoder.decode(record.data);
            const recordType = record.recordType === "url" ? "url" : "text";
            resolve({ tagContent, recordType });
          };

          ndef.onreadingerror = (event) => {
            reject(event.error);
          };
        }),
      ]);
    } catch (e) {
      throw new Error("Error reading NFC tag: " + e.message);
    }
  }

  /**
   * Write data to an NFC tag.
   * @param {NFCTagData} data
   * @returns {Promise<void>}
   */
  async function writeNFCTag(data) {
    try {
      await Promise.race([
        new Promise((_, reject) => {
          setTimeout(() => {
            reject(new Error("operation timed out"));
          }, 5000);
        }),
        new Promise((resolve, reject) => {
          const ndef = new NDEFReader();
          ndef
            .write({
              records: [{ recordType: data.recordType, data: data.tagContent }],
            })
            .then(resolve)
            .catch(reject);
        }),
      ]);
    } catch (error) {
      throw new Error("Error writing NFC tag: " + error.message);
    }
  }

  /**
   * Make NFC tag read-only.
   * @returns {Promise<void>}
   */
  async function makeNFCReadOnly() {
    try {
      await Promise.race([
        new Promise((_, reject) => {
          setTimeout(() => {
            reject(new Error("operation timed out"));
          }, 3000);
        }),
        new Promise((resolve, reject) => {
          if (
            !("NDEFReader" in window && "makeReadOnly" in NDEFReader.prototype)
          ) {
            reject(
              new Error(
                "your browser doesn't support setting an NFC tag to read only"
              )
            );
          }

          const ndef = new NDEFReader();
          ndef.makeReadOnly().then(resolve).catch(reject);
        }),
      ]);
    } catch (e) {
      throw new Error("Error making NFC tag read only: " + e.message);
    }
  }

  /**
   * Class for continuously reading NFC tags.
   * @class
   */
  class NFCContinuousRead {
    /**
     * @param {{
     *   onTagContent: (tagContent: string) => void,
     *   onError: (error: Error) => void
     * }} callbacks
     */
    constructor({ onTagContent, onError }) {
      this.ndef = new NDEFReader();
      this.onTagContent = onTagContent;
      this.onError = onError;
    }

    /**
     * Start continuous NFC tag reading.
     */
    start() {
      this.ndef
        .scan()
        .then(() => {
          this.ndef.onreading = (event) => {
            const record = event.message.records[0];
            const decoder = new TextDecoder(record.encoding);
            const tagContent = decoder.decode(record.data);
            this.onTagContent(tagContent);
          };
          this.ndef.onreadingerror = (event) => {
            this.onError(event.error);
          };
        })
        .catch((error) => {
          this.onError(error);
        });
    }
  }

  return {
    checkDeviceHasNFC,
    readNFCTag,
    writeNFCTag,
    makeNFCReadOnly,
    NFCContinuousRead,
  };
})();
