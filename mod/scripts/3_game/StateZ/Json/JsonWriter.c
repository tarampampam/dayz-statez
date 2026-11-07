/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Utility for writing JSON files atomically.
    \details Writes JSON to a .tmp file, then copies it to the target .json file and removes the temp,
             so consumers never read a partially written file.
*/
class StateZ_JsonWriter
{
	/**
	    \brief Write a JSON string atomically to the data directory as name.json.
	    \param name File name without extension.
	    \param json JSON content to write.
	    \return true on success.
	*/
	static bool Write(string name, string json)
	{
		if (!StateZ_Config.IsLoaded())
			return false;

		string target = StateZ_Constants.DATA_DIR + name + ".json";
		string tmp = StateZ_Constants.DATA_DIR + name + ".tmp";

		if (!WriteFile(tmp, json))
			return false;

		DeleteFile(target);

		if (!CopyFile(tmp, target)) {
			ErrorEx(
			    string.Format("StateZ: atomic write failed '%1' -> '%2'", tmp, target),
			    ErrorExSeverity.ERROR);
			DeleteFile(tmp);
			return false;
		}

		DeleteFile(tmp);
		return true;
	}

	/**
	    \brief Write string to file, overwriting existing content.
	    \param path Full file path.
	    \param content Content to write.
	    \return true on success.
	*/
	protected static bool WriteFile(string path, string content)
	{
		FileHandle fh = OpenFile(path, FileMode.WRITE);
		if (!fh) {
			ErrorEx(
			    string.Format("StateZ: open file '%1' failed", path),
			    ErrorExSeverity.ERROR);
			return false;
		}

		FPrint(fh, content);
		CloseFile(fh);
		return true;
	}

	/**
	    \brief Delete all .json and .tmp files in the data directory.
	*/
	static void DeleteAll()
	{
		DeleteByPattern(StateZ_Constants.DATA_DIR, "*.json");
		DeleteByPattern(StateZ_Constants.DATA_DIR, "*.tmp");
	}

	/**
	    \brief Delete all files matching a pattern in a directory.
	    \param dir Target directory (must end with /).
	    \param pattern Glob pattern (e.g. "*.json").
	*/
	protected static void DeleteByPattern(string dir, string pattern)
	{
		string fileName;
		FileAttr attr;
		FindFileHandle fh = FindFile(dir + pattern, fileName, attr, FindFileFlags.ARCHIVES);
		if (!fh)
			return;

		if (fileName != string.Empty)
			DeleteFile(dir + fileName);

		while (FindNextFile(fh, fileName, attr)) {
			if (fileName != string.Empty)
				DeleteFile(dir + fileName);
		}

		CloseFindFile(fh);
	}
}
#endif
